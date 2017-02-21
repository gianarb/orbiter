package core

import "testing"

func TestParseDocumentatinWithSingleAutoscaler(t *testing.T) {
	var data = `
autoscalers:
  aws-general:
    provider: aws
    parameters:
      aws_key: 4gegxrt5hxrht6ht
      aws_secret: rgxrtbxrtbrtbrt
    policies:
      micro:
        up: 2
        down: 3

      3d2b152bc3f6:
        up: 2
        down: 3`
	conf, err := ParseYAMLConfiguration([]byte(data))
	t.Log(conf)
	if err != nil {
		t.Fatal(err)
	}
	if conf.AutoscalersConf["aws-general"].Policies["micro"].Up != 2 {
		t.Fatal("micro expects Up equals 2")
	}
}

func TestParseDocumentatinWithMultipleAutoscaler(t *testing.T) {
	var data = `
autoscalers:
  swarm-first:
    provider: swarm
    policies:
      3d2b152bc3f6:
        up: 2
        down: 3
  aws-general:
    provider: aws
    parameters:
      aws_key: 4gegxrt5hxrht6ht
      aws_secret: rgxrtbxrtbrtbrt
    policies:
      3d2b152bc3f6:
        up: 2
        down: 3`
	conf, err := ParseYAMLConfiguration([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(conf.AutoscalersConf) != 2 {
		t.Fatal("micro expects Up equals 2")
	}
}
