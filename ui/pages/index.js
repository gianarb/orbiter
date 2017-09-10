import Head from 'next/head'
import Link from 'next/link'

import { Container, Header, Menu, Segment } from 'semantic-ui-react'

export default () => (
  <div>
    <Head>
      <link rel='stylesheet' href='//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.12/semantic.min.css' />
    </Head>
    <Menu fixed='top' inverted>
      <Container>
        <Link href='/'>
          <a className='header item'>
            <img className='ui mini image' src='/static/logo-round.svg' style={{ marginRight: '1.5em' }} /> Orbiter UI
          </a>
        </Link>
      </Container>
    </Menu>
    <Container style={{ marginTop: '6em' }}>
      <Header as='h1' dividing>Dashboard</Header>
      <Segment>Service 1</Segment>
      <Segment>Service 2</Segment>
      <Segment>Service 3</Segment>
      <Segment>Service 4</Segment>
      <Segment>Service 5</Segment>
    </Container>
  </div>
)
