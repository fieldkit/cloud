
import React, { PropTypes } from 'react'
import { map, constrain } from '../../utils.js'

class Timeline extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
      mouseOver: -1
    }
  }

  render () {

    const { currentDate, startDate, endDate, updateDate } = this.props
    const { mouseOver } = this.state

    const lineWidth = 0.01388888889 * window.innerWidth
    const lineHeight = 0.02487562189 * window.innerHeight
    const h = (window.innerHeight - lineHeight * 2) - lineHeight * 2
    const w = mouseOver >= 0 && mouseOver < h ? lineWidth / 2 : 4
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
            updateDate(nextDate)
          }}
        >
          <rect
            x={ x }
            y={ y }
            width={ w }
            height={ h }
            fill="rgba(255,255,255,0.25)"
          />
          <rect
            x={ x }
            y={ y }
            width={ w }
            height={ h }
            fill="rgba(255,255,255,0.5)"
          />
          <rect
            x={ x }
            y={ y }
            width={ w }
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
              x={ x }
              style={{
                pointerEvents: 'none'
              }}
              y={ mouseOver + lineHeight - 2 }
              width={ w }
              height={ 2 }
              fill="rgba(255,255,255,0.75)"
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



