import React from 'react'

import 'isomorphic-fetch' /* global fetch */

import {Segment} from 'semantic-ui-react'

import Layout from '../components/layout'

export default class extends React.Component {
  static async getInitialProps ({ query }) {
    const res = await fetch(`${process.env.API_HOST}/handle/${query.name}/inspect-service`)
    const json = await res.json()
    return {
      service: json
    }
  }

  render () {
    return (
      <Layout title='Service info'>
        <Segment>
          <pre style={{ margin: '0' }}>{ JSON.stringify(this.props.service, null, '  ') }</pre>
        </Segment>
      </Layout>
    )
  }
}
