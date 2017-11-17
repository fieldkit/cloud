import React, { PropTypes } from 'react'
import { connect } from 'react-redux'

class MapPage extends React.Component {
    constructor() {
        super()
        console.log("MapPage")
    }

    render() {
        return (
            <div>
                Map
            </div>
        )
    }
}

function mapStateToProps(state) {
    return {
    };
}

export default connect(mapStateToProps, {
})(MapPage);
