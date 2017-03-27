
import React from 'react'
import { map } from '../../../utils.js'

class ActivityGraph extends React.Component {
  render () {
    const {
      documents,
      startDate,
      endDate,
      x,
      y,
      w,
      h,
      rad
    } = this.props

    const lineWidth = 1.388888889
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
            const width = map(b, 0, maxDocuments, w, w + (lineWidth * window.innerWidth / 100) * 2)
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