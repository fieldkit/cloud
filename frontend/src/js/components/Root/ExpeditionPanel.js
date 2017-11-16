import React from 'react'

import logoFieldKit from '../../../img/fieldkit-logo-red.svg'
import iconClose from '../../../img/icon-close.svg'

class ExpeditionPanel extends React.Component {
    constructor(props) {
        super(props)
        this.state = {}
        this.stopPropagation = this.stopPropagation.bind(this)
    }

    shouldComponentUpdate(props) {
        return this.props.expeditionPanelOpen !== props.expeditionPanelOpen
    }

    stopPropagation(e) {
        e.stopPropagation()
    }

    render() {
        const { expeditionPanelOpen, project, closeExpeditionPanel } = this.props

        return (
            <div className={ `expedition-panel ${expeditionPanelOpen ? 'active' : ''}` } onClick={closeExpeditionPanel}>
                <div className="expedition-panel_content" onClick={this.stopPropagation}>
                    <div className="expedition-panel_content_header">
                        <div className="expedition-panel_content_header_close" onClick={closeExpeditionPanel}>
                            <img src={'/' + iconClose} />
                        </div>
                        <h1 className="expedition-panel_content_header_title">
                          {project.get('name')}
                        </h1>
                    </div>
                    <div className="expedition-panel_content_body">
                    </div>
                </div>
            </div>
        )
    }
}

export default ExpeditionPanel
