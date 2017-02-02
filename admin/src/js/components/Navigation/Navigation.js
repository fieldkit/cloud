import React, { PropTypes } from 'react'
// import NavigationItem from './NavigationItem'
import { Link } from 'react-router'

import fieldkitLogo from '../../../img/fieldkit-logo.svg'
import profileButton from '../../../img/profile-button.png'
import backgroundImage from '../../../img/bkg.jpg'
import linkImage from '../../../img/link.svg'

class Navigation extends React.Component {

  render () {
    const { 
      projectID,
      expeditionID,
      expeditions,
      projects,
      requestSignOut
    } = this.props

    return (
      <div id="side-bar">
        <div
          className="side-bar_background"
          style={{
            backgroundImage: 'url(\'/' + backgroundImage + '\')'
          }}
        />
        <div className="side-bar_logo">
          <Link to={'/admin'}>
            <img src={'/' + fieldkitLogo} alt="fieldkit logo" />
          </Link>
        </div>

        <div className="side-bar_navigation">
          <ul className="side-bar_navigation_profile">
            <Link to={'/admin'}>
              <li>Ian Ardouin-Fumat</li>
            </Link>
            <a
              href="#"
              onClick={(e) => {
                e.preventDefault()
                requestSignOut()
              }}
            >
              <li>x</li>
            </a>
          </ul>
          <div
            className="side-bar_navigation_slider"
            style={{
              marginLeft: !projectID ? 0 : '-100%'
            }}
          >
            <ul className="side-bar_navigation_slider_projects">
              <li className="side-bar_navigation_slider_sub-navigation">
                {
                  !!projects &&
                  <Link to={'/admin/new-project'}>
                    <li className="side-bar_navigation_slider_sub-navigation_button">
                      New Project <span>+</span>
                    </li>
                  </Link>
                }
              </li>
              {
                !!projects &&
                projects.map(project => {
                  return (
                    <Link to={'/admin/' + project.get('id')}>
                      <li className="side-bar_navigation_slider_projects_project">
                        <h4>
                          { project.get('name') }
                        </h4>
                      </li>
                    </Link>
                  )
                })
              }
              {
                !projects &&
                <li className="side-bar_navigation_slider_spinning-wheel"></li>
              }
            </ul>
            {
              !!projects && !!projectID && !!expeditions &&
              <ul className="side-bar_navigation_slider_expeditions">
                <li className="side-bar_navigation_slider_sub-navigation">
                  <Link to={'/admin'}>
                    <div className="side-bar_navigation_slider_sub-navigation_button">
                      <span></span>
                      Projects
                    </div>
                  </Link>
                  <Link to={'/admin/' + projectID + '/new-expedition'}>
                    <div className="side-bar_navigation_slider_sub-navigation_button">
                      New Expedition +
                    </div>
                  </Link>
                </li>
                {
                  expeditions.map(expedition => {
                    return (
                      <li className="side-bar_navigation_slider_expeditions_expedition">
                        <Link to={'/admin/' + projectID + '/' + expedition.get('id')}>
                          <h4>
                            { expedition.get('name') } 
                            <a href={'https://' + projectID + '.fieldkit.org/' + expedition.get('id')}>
                              <span
                                style={{
                                  backgroundImage: 'url(\'/' + linkImage + '\')'
                                }}
                              ></span>
                            </a>
                          </h4>
                        </Link>
                        { 
                          expedition.get('id') === expeditionID &&
                          <ul className="sections">
                            <li>
                              <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/general-settings'}>
                                General Settings
                              </Link>
                            </li>
                            <li>
                              <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/inputs'}>
                                Data inputs
                              </Link>
                            </li>
                          </ul>
                        }
                      </li>
                    )
                  })
                }
              </ul>
            }
            {
              !!projects && !!projectID && !expeditions &&
              <ul className="side-bar_navigation_slider_expeditions">
                <Link to={'/admin'}>
                  <li className="side-bar_navigation_slider_expeditions_back">
                    {'< Projects'}
                  </li>
                </Link>
                <li className="side-bar_navigation_slider_expeditions_title">
                  { projects.getIn([projectID, 'name']) }
                </li>
                <li className="side-bar_navigation_slider_spinning-wheel"></li>
              </ul>
            }
          </div>
        </div>
      </div>
    )
  }
}

Navigation.propTypes = {
}

export default Navigation