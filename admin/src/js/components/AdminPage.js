import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import NavigationPanel from './NavigationPanel'
import BreadCrumbs from './BreadCrumbs'
import Modal from './Modal'

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

    const { 
      children, 
      params, 
      disconnect,
      location,
      modal,
      cancelAction,
      saveChangesAndResume
    } = this.props

    const { expeditions } = this.state

    const modalProps = {
      modal,
      cancelAction,
      saveChangesAndResume
    }

    const modalComponent = () => {
      if (!!modal.get('type')) {
        return (
          <Modal {...modalProps} />
        )
      } else {
        return null
      }
    }

    return (
      <div id="admin-page" className="page">
        { modalComponent() }
        <NavigationPanel {...params} expeditions={expeditions} disconnect={disconnect} />
        <div className="page-content">
          <BreadCrumbs {...location} />
          {children}
        </div>
      </div>
    )
  }
}

AdminPage.propTypes = {
  disconnect: PropTypes.func.isRequired,
  saveChangesToTeam: PropTypes.func.isRequired,
  clearChangesToTeam: PropTypes.func.isRequired
}

export default AdminPage