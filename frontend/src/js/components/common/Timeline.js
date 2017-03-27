
import React, { PropTypes } from 'react'
import { map, constrain } from '../../utils.js'
import ActivityGraph from './Timeline/ActivityGraph'

class Timeline extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
      mouseOver: -1
    }
  }

  render () {
    const {
      currentDate,
      startDate,
      endDate,
      documents,
      updateDate
    } = this.props
    const { mouseOver } = this.state

    if (!documents || documents.size === 0) return null

    const lineWidth = 0.01388888889 * window.innerWidth
    const lineHeight = 0.02487562189 * window.innerHeight
    const h = (window.innerHeight - lineHeight * 2) - lineHeight * 2
    const w = 4
    const hoverWidth = mouseOver >= 0 && mouseOver < h ? (lineWidth / 2 - w) : 0
    const rad = lineWidth / 2
    const x = lineWidth - w / 2
    const y = lineHeight
    const progress = map(currentDate, startDate, endDate, 0, 1)

    return (
      <div id="timeline">
        <svg
          onMouseMove={(e) => {
            const mouseY = e.nativeEvent.clientY - lineHeight * 2
            this.setState({
              ...this.state,
              mouseOver: mouseY
            })
          }}
          onMouseOut={() => {
            this.setState({
              ...this.state,
              mouseOver: -1
            })
          }}
          onClick={(e) => {
            const mouseY = e.nativeEvent.clientY - lineHeight * 2
            const nextDate = constrain(map(mouseY, 0, h, startDate, endDate), startDate, endDate)
            updateDate(nextDate, null, true)
          }}
        >
          <ActivityGraph
            documents={ documents }
            startDate={ startDate }
            endDate={ endDate }
            lineWidth={ lineWidth }
            x={ x }
            y={ y }
            w={ w }
            h={ h }
            rad={ rad }
          />
          <rect
            x={ x - hoverWidth / 2 }
            y={ y }
            width={ w + hoverWidth }
            height={ h }
            fill="rgb(150,150,150)"
          />
          <rect
            x={ x - hoverWidth / 2 }
            y={ y }
            width={ w + hoverWidth }
            height={ h * progress }
            fill="#D0462C"
          />
          { 
            mouseOver >= 0 && mouseOver < h &&
            <circle
              cx={ x + w / 2 }
              cy={ y + h * progress }
              r={ rad }
              fill="#D0462C"
            />
          }
          { 
            mouseOver >= 0 && mouseOver < h &&
            <rect
              x={ x - hoverWidth / 2 }
              style={{
                pointerEvents: 'none'
              }}
              y={ mouseOver + lineHeight - 2 }
              width={ w + hoverWidth }
              height={ 2 }
              fill="white"
            />
          }
        </svg>
      </div>
    )
  }

}

Timeline.propTypes = {

}

export default Timeline



