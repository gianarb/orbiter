import React from 'react'
import Link from 'next/link'

import 'isomorphic-fetch' /* global fetch */

import {Button, Grid, Header, Segment} from 'semantic-ui-react'

import Layout from '../components/layout'

const MIN_REPLICAS = 1

export default class extends React.Component {
  static async getInitialProps () {
    const res = await fetch(`${process.env.API_HOST}/autoscaler`)
    const json = await res.json()
    return {
      services: json.data
    }
  }

  constructor (props) {
    super(props)
    this.state = {
      services: props.services
    }
  }

  scaleUp (serviceName) {
    fetch(`${process.env.API_HOST}/handle/${serviceName}/up`, {
      method: 'POST',
      body: {}  // just pass the instance
    }).then(response => {
      if (response.status === 201) {
        // TODO: update react state manipulation with a better approach
        const service = this.state.services.find(item => item.name === serviceName)
        const others = this.state.services.filter(item => item.name !== serviceName)
        this.setState({ services: [...others, { name: service.name, replicas: service.replicas + 1 }] })
      }
    })
  }

  scaleDown (serviceName) {
    fetch(`${process.env.API_HOST}/handle/${serviceName}/down`, {
      method: 'POST',
      body: {}  // just pass the instance
    }).then(response => {
      if (response.status === 201) {
        // TODO: update react state manipulation with a better approach
        const service = this.state.services.find(item => item.name === serviceName)
        const others = this.state.services.filter(item => item.name !== serviceName)
        this.setState({ services: [...others, { name: service.name, replicas: service.replicas - 1 }] })
      }
    })
  }

  render () {
    return (
      <Layout title='Dashboard'>
        {
          this.state.services
            .sort((a, b) => a.name.localeCompare(b.name))
            .map(service => (
              <Segment key={service.name}>
                <Grid>
                  <Grid.Column width={8} style={{ display: 'flex', alignItems: 'center' }}>
                    <Header as='h3'>
                      <Link href={{ pathname: '/service', query: { name: service.name } }}>
                        <a>{ service.name }</a>
                      </Link>
                    </Header>
                  </Grid.Column>
                  <Grid.Column width={4} style={{ display: 'flex', alignItems: 'center' }}>
                    <Header as='h3'>{ service.replicas }</Header>
                  </Grid.Column>
                  <Grid.Column floated='right' width={4} style={{ textAlign: 'right' }}>
                    <Button icon='plus' onClick={() => { this.scaleUp(service.name) }} />
                    <Button icon='minus' disabled={service.replicas <= MIN_REPLICAS} onClick={() => { this.scaleDown(service.name) }} />
                  </Grid.Column>
                </Grid>
              </Segment>
            ))
        }
      </Layout>
    )
  }
}
