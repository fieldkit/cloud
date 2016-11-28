import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import NavigationPanel from './NavigationPanel'

class AdminPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      expeditions: [
        {
          id: 'okavango_16',
          name: 'Okavango 2016',
          startDate: new Date('2016-08-17 00:00:00+02:00')
        },
        {
          id: 'bike_angola_16',
          name: 'Bike Angola 16',
          startDate: new Date('2016-07-06 00:00:00+02:00')
        },
        {
          id: 'cuando_16',
          name: 'Cuando 16',
          startDate: new Date('2016-10-01 00:00:00+02:00')
        }
      ]
    }
  }

  render () {

    const { children, params } = this.props
    const { expeditions } = this.state

    return (
      <div id="admin-page" className="page">
        <NavigationPanel {...params} expeditions={expeditions} />
        <h1>Admin Page</h1>
        {children}
      </div>
    )
  }
}

AdminPage.propTypes = {

}

export default AdminPage