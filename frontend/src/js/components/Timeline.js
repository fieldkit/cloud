
import React, {PropTypes} from 'react'
import * as d3 from 'd3'
import autobind from 'autobind-decorator'
import SightingGraph from './SightingGraph'

class Timeline extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      padding: [2, 8],
      x: 0,
      cursorY: 0,
      width: 0,
      height: 0,
      data: [],
      range: [0, 0],
      domain: [0, 0],
      scaleDays: null,
      scaleTime: null,
      totalSightings: [],
      mouseOver: false
    }
  }

  componentWillReceiveProps (nextProps) {
    const { expedition } = nextProps
    if (expedition) {
      const { padding } = this.state
      const startDate = expedition.start
      const height = window.innerHeight - 72
      const width = window.innerWidth * 0.05
      const x = width * 0.58
      var dayCount = expedition.dayCount + 1

      var data = []
      for (var i = 0; i < dayCount; i++) {
        var d = new Date(startDate.getTime() + i * (1000 * 3600 * 24))
        data.push(d)
      }

      const range = [0 + padding[1], height - padding[1]]
      const domain = [dayCount - 1, 0]
      const scaleDays = d3.scaleLinear()
        .domain(domain)
        .range(range)
      const scaleTime = d3.scaleLinear()
        .domain([startDate.getTime() + (dayCount - 1) * (1000 * 3600 * 24), startDate.getTime()])
        .range(range)

      const cursorY = this.state.mouseOver ? this.state.cursorY : (scaleTime(expedition.currentDate.getTime()) - 8)

      var totalSightings = expedition.totalSightings

      this.setState({
        ...this.state,
        x,
        cursorY,
        width,
        height,
        data,
        range,
        domain,
        scaleDays,
        scaleTime,
        totalSightings
      })
    }
  }

  @autobind
  onClick (e) {
    const { scaleTime, range } = this.state
    const { jumpTo, expeditionID } = this.props
    var y = e.nativeEvent.offsetY
    jumpTo(new Date(scaleTime.invert(Math.max(range[0] + 1, Math.min(range[1] - 1, y)))), expeditionID)
  }

  @autobind
  onMouseMove (e) {
    const { range, padding } = this.state
    const cursorY = Math.max(range[0], Math.min(range[1], e.nativeEvent.offsetY)) - padding[1]
    this.setState({
      ...this.state,
      cursorY,
      mouseOver: true
    })
  }

  @autobind
  onMouseOut (e) {
    const { expedition } = this.props
    const { scaleTime } = this.state
    const cursorY = (scaleTime(expedition.currentDate.getTime()) - 8)
    this.setState({
      ...this.state,
      cursorY,
      mouseOver: false
    })
  }

  render () {
    const { expedition } = this.props
    if (!expedition) return <svg id="timeline"></svg>
    const { width, height, data, x, range, scaleDays, cursorY, totalSightings, padding } = this.state

    const days = data.map((d, i) => {
      return <circle cx={x} cy={scaleDays(i)} r={3} key={i} fill="white"/>
    })

    return (
      <svg
        id="timeline"
        className={location.pathname.indexOf('/about') > -1 || location.pathname.indexOf('/data') > -1 ? 'invisible' : 'visible'}
        style={{height: height + 'px'}}
        onMouseOut={this.onMouseOut}
        onMouseMove={this.onMouseMove}
        onClick={this.onClick}
      >
        <filter id="dropshadow" height="120%">
          <feGaussianBlur in="SourceAlpha" stdDeviation="3"/>
          <feOffset dx="2" dy="0" result="offsetblur"/>
          <feMerge>
            <feMergeNode/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
        <g transform={'translate(' + 0 + ',' + padding[1] + ')'} style={{pointerEvents: 'none'}}>
          <SightingGraph sightings={totalSightings} width={width} height={height - padding[1] * 2}/>
        </g>
        <line x1={x} x2={x} y1={range[0]} y2={range[1]} style={{stroke: 'white'}}/>
        <g>{ days }</g>
        <g transform={'translate(' + (x - 20) + ',' + cursorY + ')'} style={{pointerEvents: 'none'}}>
          <path fill="#F9D144" d="M8,0c5,0,12,8,12,8s-7,8-12,8c-4.4,0-8-3.6-8-8C0,3.6,3.6,0,8,0z" style={{filter: 'url(#dropshadow)'}}/>
          <circle fill="#1F1426" cx="7.9" cy="7.8" r="3"/>
        </g>
      </svg>
    )
  }
}

Timeline.propTypes = {
  expedition: PropTypes.object,
  jumpTo: PropTypes.func.isRequired,
  expeditionID: PropTypes.string
}
export default Timeline
