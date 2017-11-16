
import React, { PropTypes } from 'react'
import { map, constrain } from '../../../utils.js'
import ActivityGraph from './ActivityGraph'
import { getLineHeight, getLineWidth } from '../../../utils.js'

class Timeline extends React.Component {

    constructor(props) {
        super(props)
        this.state = {
            mouseOver: -1
        }
        this.onMouseMove = this.onMouseMove.bind(this)
        this.onMouseOut = this.onMouseOut.bind(this)
        this.onClick = this.onClick.bind(this)
    }

    onMouseMove(e) {
        const lineHeight = getLineHeight()
        const mouseY = e.nativeEvent.clientY - lineHeight * 2
        this.setState({
            ...this.state,
            mouseOver: mouseY
        })
    }

    onMouseOut(e) {
        this.setState({
            ...this.state,
            mouseOver: -1
        })
    }

    onClick(e) {
        const lineHeight = getLineHeight()
        const h = (window.innerHeight - lineHeight * 2) - lineHeight * 2
        const {updateDate, startDate, endDate} = this.props
        const mouseY = e.nativeEvent.clientY - lineHeight * 2
        const nextDate = constrain(map(mouseY, 0, h, startDate, endDate), startDate, endDate)
        updateDate(nextDate, null, true)
    }

    render() {
        const {currentDate, startDate, endDate, documents, } = this.props
        const {mouseOver} = this.state

        if (!documents || documents.size === 0) return null

        const lineWidth = getLineWidth()
        const lineHeight = getLineHeight()
        const h = (window.innerHeight - lineHeight * 2) - lineHeight * 2
        const w = 4
        const hoverWidth = mouseOver >= 0 && mouseOver < h ? (lineWidth / 2 - w) : 0
        const rad = lineWidth / 2
        const x = lineWidth - w / 2
        const y = lineHeight
        const progress = map(isNaN(currentDate) ? startDate : currentDate, startDate, endDate, 0, 1)

        return (
            <div id="timeline">
                <svg onMouseMove={ this.onMouseMove } onMouseOut={ this.onMouseOut } onClick={ this.onClick }>
                    <ActivityGraph documents={ documents } startDate={ startDate } endDate={ endDate } lineWidth={ lineWidth } x={ x } y={ y }
                        w={ w } h={ h } rad={ rad } />
                    <rect x={ x - hoverWidth / 2 } y={ y } width={ w + hoverWidth } height={ h } fill="rgb(150,150,150)" />
                    <rect x={ x - hoverWidth / 2 } y={ y } width={ w + hoverWidth } height={ h * progress } fill="#D0462C" />
                    { mouseOver >= 0 && mouseOver < h &&
                      <circle cx={ x + w / 2 } cy={ y + h * progress } r={ rad } fill="#D0462C" /> }
                    { mouseOver >= 0 && mouseOver < h &&
                      <rect x={ x - hoverWidth / 2 } style={ { pointerEvents: 'none' } } y={ mouseOver + lineHeight - 2 } width={ w + hoverWidth } height={ 2 } fill="white" /> }
                </svg>
            </div>
        )
    }

}

Timeline.propTypes = {

}

export default Timeline



