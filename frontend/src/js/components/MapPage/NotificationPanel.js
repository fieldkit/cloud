import React, { PropTypes } from 'react'
import { dateToString } from '../../utils.js'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import { is } from 'immutable'

import DocumentCard from '../../components/common/DocumentCard'

class NotificationPanel extends React.Component {
    shouldComponentUpdate(props) {
        return !is(this.props.currentDocuments, props.currentDocuments)
    }

    render() {
        const {currentDocuments} = this.props

        const panels = Array.from(currentDocuments.values())
              .map((d, i) => <DocumentCard document={d} key={i} className={"notification-panel_post"} />)

        return (
            <div className="notification-panel">
                <ReactCSSTransitionGroup transitionName="transition" transitionEnter={true} transitionEnterTimeout={500} transitionLeave={true}
                    transitionLeaveTimeout={500}>
                    {panels}
                </ReactCSSTransitionGroup>
            </div>
        )
    }
}

NotificationPanel.propTypes = {
}

export default NotificationPanel
