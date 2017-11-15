import React, { PropTypes } from 'react'
import { is } from 'immutable'

import logo from "../../../img/fieldkit-logo.svg"

class DisclaimerPanel extends React.Component {

    render() {
        const {} = this.props

        return (
            <div className="disclaimer-panel">
                <div className="disclaimer-body">
                    <span className="b">NOTE:</span> Map images have been obtained from a third-party and do not reflect the editorial decisions of National Geographic.
                </div>
                <img className="disclaimer-logo" src={ logo } />
            </div>
        )
    }

}

DisclaimerPanel.propTypes = {

}

export default DisclaimerPanel
