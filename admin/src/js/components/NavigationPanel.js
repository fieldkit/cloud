import React, { PropTypes } from 'react'
import NavigationItem from './NavigationItem'
import { Link } from 'react-router'

class NavigationPanel extends React.Component {

  render () {
    const { projectID, expeditionID, expeditions, disconnect, projects } = this.props

    return (
      <div id="header">
        <div className="background"/>
        <div id="logo">
          <Link to={'/admin'}>
            <img src="/src/img/fieldkit-logo.svg" alt="fieldkit logo" />
          </Link>
          <Link to={'/admin/profile'}>
            <img src="/src/img/profile-button.png" />
          </Link> 
          {/*
          <SignOutButton
            endpoint={'localhost:3000/signout'}
            next={disconnect}
          >
            X
          </SignOutButton>
          */}
        </div>
        {
          projects
            .map(project => {
              return (
                <div className="project">
                  <div className="project-name">
                    <h3>
                      { project.get('name') }
                    </h3>
                  </div>
                  <ul className="expeditions">
                    {
                      expeditions
                        .filter(e => {
                          const projectExpeditions = project.get('expeditions')
                          return projectExpeditions.includes(e.get('id'))
                        })
                        .map(expedition => {
                          return (
                            <NavigationItem 
                              expedition={expedition}
                              active={expeditionID === expedition.get('id')}
                              key={expedition.get('id')} 
                              projectID={projectID}
                            />
                          )
                        })
                    }
                    <li className="new-expedition">
                      <Link to={'/admin/' + projectID + '/new-expedition'}>
                        <h4>Add New Expedition</h4>
                      </Link>
                    </li>
                  </ul>
                </div>
              )
            })
        }
        <div className="project">
          <div className="project-name">
            <h3>
              Add New Project
            </h3>
          </div>
        </div>
      </div>
    )
  }
}

NavigationPanel.propTypes = {
  disconnect: PropTypes.func.isRequired
}

export default NavigationPanel