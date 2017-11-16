import React from 'react'

import sensorIcon from '../../../img/icon-sensor-red.svg'
import twitterIcon from '../../../img/icon-twitter.png'

export default class DocumentCard extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        const { document: d, style, className } = this.props

        let title, body, icon, extra;
        if (d.get("user")) {
            title = `@${d.get("user").get("screen_name")}`
            body = d.get("text")
            icon = `/${twitterIcon}`
            if (d.getIn(["entities", "media"])) {
                let url = d.getIn(["entities", "media", 0, "media_url_https"])
                extra = (<img className='notification-panel_twitter-media' src={ url } />)
            } else {
                extra = (<span></span>)
            }
        } else if (d.get("GPSSpeed")) {
            title = `Sensor (${d.get("SampleType")})`
            body = `GPS Speed: ${d.get("GPSSpeed")}`
            icon = `/${sensorIcon}`
            extra = (<span>TODO</span>)
        } else {
            title = `Sensor`
            body = "Detailed sensor data will go here."
            icon = `/${sensorIcon}`
            extra = (<span>TODO</span>)
        }

        return (
            <div className={className + " doc-card"} style={style}>
                <div className={className + "_content"}>
                    <div className={className + "_content_icon"}>
                        <img src={icon} width="100%" />
                    </div>
                    <div className={className + "_title"}> {title} </div>
                    <div>
                        {body}
                    </div>
                    {extra}
                </div>
            </div>
        )
    }
}
