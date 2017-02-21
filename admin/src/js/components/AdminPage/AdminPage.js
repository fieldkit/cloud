import React, {PropTypes, Children} from 'react'
import { Link } from 'react-router'
import Navigation from '../Navigation/Navigation'
import BreadCrumbs from './BreadCrumbs'
import Modal from './Modal'

class AdminPage extends React.Component {

  render () {

    const { 
      params,
      errors,
      location,
      modal,
      cancelAction,
      saveChangesAndResume,
      expeditions,
      projects,
      breadcrumbs,
      user,
      requestSignOut
    } = this.props

    const modalProps = {
      modal,
      cancelAction,
      saveChangesAndResume
    }

    const children = Children.map(this.props.children,
      (child) => React.cloneElement(child, {
        errors
      })
    )

    const modalComponent = () => {
      if (!!modal.get('type')) {
        return (
          <Modal { ...modalProps } />
        )
      } else {
        return null
      }
    }

    return (
      <div id="admin-page" className="page">
        { modalComponent() }
        <Navigation 
          { ...params }
          expeditions={ expeditions }
          projects={ projects }
          requestSignOut={ requestSignOut }
          user={ user }
        />
        <div className="page-content">
          <BreadCrumbs 
            { ...location }
            breadcrumbs={ breadcrumbs }
          />
          { children }
        </div>
      </div>
    )
  }
}

AdminPage.propTypes = {
}

export default AdminPage