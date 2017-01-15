
import React, { PropTypes } from 'react'
import { map } from '../utils.js'

class Timeline extends React.Component {

  constructor (props) {
    super(props)
    this.sate = {}
  }

  render () {

    const { currentDate, startDate, endDate } = this.props

    const lineWidth = 0.01388888889 * window.innerWidth
    const lineHeight = 0.02487562189 * window.innerHeight
    const w = 5
    const h = (window.innerHeight - lineHeight * 2) - lineHeight * 2
    const rad = Math.min(lineWidth, lineHeight) / 2
    const x = lineWidth - w / 2
    const y = lineHeight
    const progress = map(currentDate, startDate, endDate, 0, 1)

    return (
      <div id="timeline">
        <svg>
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
          <circle
            cx={ x + w / 2 }
            cy={ y + h * progress }
            r={ rad }
            fill="#D0462C"
          />
        </svg>
      </div>
    )
  }

}

Timeline.propTypes = {

}

export default Timeline



