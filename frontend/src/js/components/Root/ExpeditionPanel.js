
import React from 'react'

import logoFieldkit from '../../../img/fieldkit-logo-red.svg'
import logoNatGeo from '../../../img/national-geographic-logo-long.svg'
import iconClose from '../../../img/icon-close.svg'

class ExpeditionPanel extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
    this.stopPropagation = this.stopPropagation.bind(this)
  }

  shouldComponentUpdate (props) {
    return this.props.expeditionPanelOpen !== props.expeditionPanelOpen
  }

  stopPropagation (e) {
    e.stopPropagation()
  }

  render () {
    const {
      expeditionPanelOpen,
      project,
      expeditions,
      currentExpedition,
      closeExpeditionPanel
    } = this.props

    let hostname = location.hostname.split('.').slice(1).join('.')
    const protocol = hostname !== 'localhost' ? 'https://' : ''
    if (hostname === 'localhost') hostname += ':' + location.port

    return (
      <div
        className={`expedition-panel ${ expeditionPanelOpen ? 'active' : '' }`}
        onClick={ closeExpeditionPanel }
      >
        <div
          className="expedition-panel_content"
          onClick={ this.stopPropagation }
        >
          <div className="expedition-panel_content_header">
            <div
              className="expedition-panel_content_header_close"
              onClick={ closeExpeditionPanel }
            >
              <img src={ '/' + iconClose }/>
            </div>
            <h1 className="expedition-panel_content_header_title">
              { project.get('name') }
            </h1>
          </div>
          <div className="expedition-panel_content_body">
            <p className="expedition-panel_content_body_label">
              Other expeditions from the same project:
            </p>
            <ul className="expedition-panel_content_body_projects">
              {
                expeditions
                  .filter(e => {
                    return e.get('id') !== currentExpedition
                  })
                  .map(e => {
                    return (
                      <a href={ `${protocol}${project.get('id')}.${hostname}/${e.get('id')}` }>
                        <li className="expedition-panel_content_body_projects_project">
                          <h2>
                            { e.get('name') }
                          </h2>
                        </li>
                      </a>
                    )
                  })
              }
            </ul>
            <div className="expedition-panel_content_body_separator"/>

            <a href="https://fieldkit.org">
              <div className="expedition-panel_content_body_logos">
                <img src={ '/' + logoNatGeo } />
                <img src={ '/' + logoFieldkit } />
              </div>
            </a>
            <p>
              Fieldkit is a one-click open platform for field researchers and explorers.
            </p>
          </div>
        </div>
      </div>
    )
  }
}

export default ExpeditionPanel