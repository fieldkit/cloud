import React, { PropTypes, Component } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import DocumentCard from '../../components/common/DocumentCard'

export default class WebGLOverlay extends Component {
    constructor(props) {
        super(props)
        this.state = {}
        this.onClick = this.onClick.bind(this)
    }

    onClick() {
        const { focusedDocument, openLightbox } = this.props
        if (!!focusedDocument) {
            openLightbox(focusedDocument.get('id'))
        }
    }

    render() {
        const { focusedDocument } = this.props

        let card = null

        if (focusedDocument != null) {
            card = <DocumentCard document={focusedDocument} className={"dom-overlay_popup"} style={{ left: focusedDocument.get('x'), top: focusedDocument.get('y') }} />
        }

        return (
            <div className="dom-overlay" style={{ cursor: !!focusedDocument ? 'pointer' : 'auto' }} onClick={this.onClick}>
                <ReactCSSTransitionGroup transitionName="transition" transitionEnter={true} transitionEnterTimeout={200} transitionLeave={true} transitionLeaveTimeout={200}>
                    {card}
                </ReactCSSTransitionGroup>
            </div>
        )
    }
}

WebGLOverlay.propTypes = {
}
