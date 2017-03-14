import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class DashboardSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  render () {

    const {
      currentExpedition
    } = this.props

    if (!currentExpedition) {
      return (
        <li className="spinning-wheel"></li>
      )
    }

    return (
      <div id="dashboard-section" className="section">
        <h1>{currentExpedition.get('name')}</h1>
        <h2>Dashboard section</h2>
      </div>
    )
  }
}

DashboardSection.propTypes = {

}

export default DashboardSection