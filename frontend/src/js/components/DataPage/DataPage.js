
import React, { PropTypes } from 'react'

class DataPage extends React.Component {
    render() {
        return (
            <div className="data-page page">
                <div className="data-page_map-overlay" />
                <div className="data-page_content" ref="content">
                    This page will contain the API documentation and explorer...
                </div>
            </div>
        )
    }

}

DataPage.propTypes = {

}

export default DataPage
