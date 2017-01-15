
import React, { PropTypes } from 'react'

class Timeline extends React.Component {

  constructor (props) {
    super(props)
    this.sate = {}
  }

  render () {

    const lineWidth = 0.01388888889 * window.innerWidth
    const lineHeight = 0.02487562189 * window.innerHeight
    const w = 5
    const h = (window.innerHeight - lineHeight * 2) - lineHeight * 2
    const rad = Math.min(lineWidth, lineHeight) / 2
    const x = lineWidth - w / 2
    const y = lineHeight

    return (
      <div id="timeline">
        <svg>
          <rect
            x={ x }
            y={ y }
            width={ w }
            height={ h }
            fill="rgba(255,255,255,0.3)"
          />
          <rect
            x={ x }
            y={ y }
            width={ w }
            height={ h / 2 }
            fill="rgba(255,255,255,0.6)"
          />
          <rect
            x={ x }
            y={ y }
            width={ w }
            height={ h / 4 }
            fill="#D0462C"
          />
          <circle
            cx={ x + w / 2 }
            cy={ y + h / 4 }
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



