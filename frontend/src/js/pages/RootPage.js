import React, { PropTypes } from 'react'
import { connect } from 'react-redux'

class RootPage extends React.Component {
    constructor() {
        super()
        console.log("RootPage")
    }

    render() {
        return (
            <div>
                Root
            </div>
        )
    }
}

function mapStateToProps(state) {
    return {
    };
}

export default connect(mapStateToProps, {
})(RootPage);
