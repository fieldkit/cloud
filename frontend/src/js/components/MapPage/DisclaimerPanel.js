import React, { PropTypes } from 'react'

class DisclaimerPanel extends React.Component {

    render() {
        return (
            <div className="disclaimer-panel">
                <div className="disclaimer-body">
                    <span className="b">NOTE:</span> Map images have been obtained from a third-party.
                </div>
            </div>
        )
    }
}

DisclaimerPanel.propTypes = {
}

export default DisclaimerPanel
