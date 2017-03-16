package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
	"github.com/gianarb/orbiter/autoscaler"
	"golang.org/x/oauth2"
)

type DigitalOceanProvider struct {
	client *godo.Client
	config map[string]string
	ctx    context.Context
}

func NewDigitalOceanProvider(c map[string]string) (autoscaler.Provider, error) {
	tokenSource := &TokenSource{
		AccessToken: c["token"],
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	p := DigitalOceanProvider{
		client: client,
		config: c,
		ctx:    context.Background(),
	}
	return p, nil
}

func (p DigitalOceanProvider) Scale(serviceId string, target int, direction bool) error {
	var wg sync.WaitGroup
	responseChannel := make(chan response, target)

	if direction == true {
		for ii := 0; ii < target; ii++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				t := time.Now()
				i, _ := strconv.ParseInt(p.config["key_id"], 10, 64)
				createRequest := &godo.DropletCreateRequest{
					Name:     fmt.Sprintf("%s-%s", serviceId, t.Format("20060102150405")),
					Region:   p.config["region"],
					Size:     p.config["size"],
					UserData: p.config["userdata"],
					SSHKeys:  []godo.DropletCreateSSHKey{{ID: int(i)}},
					Image: godo.DropletCreateImage{
						Slug: p.config["image"],
					},
				}
				droplet, _, err := p.client.Droplets.Create(p.ctx, createRequest)
				responseChannel <- response{
					err:       err,
					droplet:   droplet,
					direction: true,
				}
			}()
		}
	} else {
		// TODO(gianarb): This can not work forever. We need to have proper pagination
		droplets, _, err := p.client.Droplets.List(p.ctx, &godo.ListOptions{
			Page:    1,
			PerPage: 500,
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"provider": "digitalocean",
				"error":    err,
			}).Warnf("Impossibile to get the list of droplets.")
			return err
		}

		ii := 0
		for _, single := range droplets {
			if p.isGoodToBeDeleted(single, serviceId) && ii < target {
				go func() {
					defer wg.Done()
					_, err := p.client.Droplets.Delete(p.ctx, single.ID)
					responseChannel <- response{
						err:       err,
						droplet:   &single,
						direction: false,
					}
				}()
				wg.Add(1)
				ii++
			}
		}

		//go func() {
		//for iii := 0; iii < target; iii++ {
		//select {
		//case err := <-errorChannel:
		//logrus.WithFields(logrus.Fields{
		//"error":    err.Error(),
		//"provider": "digitalocean",
		//}).Warnf("We was not able to delete the droplet.")
		//case droplet := <-dropletChannel:
		//logrus.WithFields(logrus.Fields{
		//"provider":    "digitalocean",
		//"dropletName": droplet.ID,
		//}).Debugf()
		//}
		//}
		//wg.Wait()
		//}()
	}
	go func() {
		var message string
		for iii := 0; iii < target; iii++ {
			r := <-responseChannel
			if r.err != nil {
				message = "We was not able to instantiate a new droplet."
				if r.direction == false {
					message = fmt.Sprintf("Impossibile to delete droplet %d ", r.droplet.ID)
				}
				logrus.WithFields(logrus.Fields{
					"error":    r.err.Error(),
					"provider": "digitalocean",
				}).Warn(message)
			} else {
				message = fmt.Sprintf("New droplet named %s with id %d created.", r.droplet.Name, r.droplet.ID)
				if r.direction == false {
					message = fmt.Sprintf("Droplet named %s with id %d deleted.", r.droplet.Name, r.droplet.ID)
				}
				logrus.WithFields(logrus.Fields{
					"provider":    "digitalocean",
					"dropletName": r.droplet.ID,
				}).Debug(message)
			}
		}
		wg.Wait()
	}()
	return nil
}

// Check if a drople is eligible to be deleted
func (p DigitalOceanProvider) isGoodToBeDeleted(droplet godo.Droplet, serviceId string) bool {
	if droplet.Status == "active" && strings.Contains(strings.ToUpper(droplet.Name), strings.ToUpper(serviceId)) {
		// TODO(gianarb): This can not work forever. We need to have proper pagination
		actions, _, _ := p.client.Droplets.Actions(p.ctx, droplet.ID, &godo.ListOptions{
			Page:    1,
			PerPage: 500,
		})
		// If there is an action in progress the droplet can not be deleted.
		for _, action := range actions {
			if action.Status == godo.ActionInProgress {
				fmt.Println(fmt.Sprintf("%d has an action in progress", droplet.ID))
				return false
			}
		}
		return true
	}
	return false
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type response struct {
	err       error
	droplet   *godo.Droplet
	direction bool
}
