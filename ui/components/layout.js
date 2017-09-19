import Link from 'next/link'
import Head from 'next/head'

import { Container, Grid, Header, List, Menu, Segment } from 'semantic-ui-react'

export default ({ children, title }) => (
  <div>
    <Head>
      <title>Orbiter UI - { title }</title>
      <link rel='stylesheet' href='//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.12/semantic.min.css' />
    </Head>
    <Segment inverted textAlign='center' vertical>
      <Container>
        <Menu fixed='top' inverted>
          <Link href='/'>
            <a className='header item'>
              <img className='ui mini image' src='/static/logo-round.svg' style={{ marginRight: '1.5em' }} /> Orbiter UI
            </a>
          </Link>
        </Menu>
      </Container>
    </Segment>
    <Segment vertical style={{ padding: '5em 0em' }}>
      <Container>
        <Header as='h1' dividing>{ title }</Header>
        { children }
      </Container>
    </Segment>
    <Segment inverted vertical style={{ padding: '3em 0em' }}>
      <Container>
        <Grid divided inverted stackable>
          <Grid.Row>
            <Grid.Column width={4}>
              <Header inverted as='h4' content='About' />
              <List link inverted>
                <List.Item as='a'>Link 1</List.Item>
                <List.Item as='a'>Link 2</List.Item>
                <List.Item as='a'>Link 3</List.Item>
                <List.Item as='a'>Link 4</List.Item>
              </List>
            </Grid.Column>
            <Grid.Column width={4}>
              <Header inverted as='h4' content='Services' />
              <List link inverted>
                <List.Item as='a'>Link 1</List.Item>
                <List.Item as='a'>Link 2</List.Item>
                <List.Item as='a'>Link 3</List.Item>
                <List.Item as='a'>Link 4</List.Item>
              </List>
            </Grid.Column>
            <Grid.Column width={8}>
              <Header as='h4' inverted>Footer Header</Header>
              <p>Description</p>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Container>
    </Segment>
  </div>
)
