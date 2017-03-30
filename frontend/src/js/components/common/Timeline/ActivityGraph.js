
import React from 'react'
import { map } from '../../../utils.js'
import { is } from 'immutable'

class ActivityGraph extends React.Component {
  shouldComponentUpdate (props) {
    return !is(this.props.documents, props.documents) ||
      this.props.lineWidth !== props.lineWidth ||
      this.props.x !== props.x ||
      this.props.y !== props.y ||
      this.props.w !== props.w ||
      this.props.h !== props.h ||
      this.props.rad !== props.rad
  }

  render () {
    const {
      documents,
      startDate,
      endDate,
      lineWidth,
      x,
      y,
      w,
      h,
      rad
    } = this.props

    const bucketCount = Math.floor(h / 5)
    const buckets = new Array(bucketCount).fill(0)
    const bucketHeight = h / bucketCount - 1

    documents.forEach(d => {
      const id = Math.floor(map(d.get('date'), startDate, endDate + 1, 0, bucketCount))
      buckets[id] ++ 
    })

    const maxDocuments = Math.max.apply(null, buckets)

    return (
      <g
        style={{
          transform: `translate(${x}px, ${y}px)`
        }}
      >
        {
          buckets.map((b, i) => {
            const width = map(b, 0, maxDocuments, w, w + lineWidth * 2)
            return (
              <rect
                key={ i }
                x={ -width / 2 + w / 2 }
                y={ map(i, 0, bucketCount, 0, h) }
                width={ width }
                height={ bucketHeight }
                fill="white"
              />
            )
          })
        }
      </g>
    )
  }
}

export default ActivityGraph