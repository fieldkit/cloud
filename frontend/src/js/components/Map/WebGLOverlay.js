import React, {PropTypes, Component} from 'react'
import React3 from 'react-three-renderer'
import ViewportMercator from 'viewport-mercator-project'
import { Vector2, Vector3, BufferAttribute, Color, VertexColors, TextureLoader, Geometry } from 'three'
import {MeshLine, MeshLineMaterial} from 'three.meshline'

import iconSensorReading from '../../../img/sighting.png'

export default class WebGLOverlay extends Component {
  constructor (props) {
    super(props)

    // here you can add other feature types and assign them a particle texture
    const bufferGeometries = {
      particles: {
        'Feature': {
          count: 1000,
          position: new BufferAttribute(new Float32Array(1000 * 3), 3),
          color: new BufferAttribute(new Float32Array(1000 * 4), 4),
          index: new BufferAttribute(new Uint16Array(1000 * 1), 1),
          texture: new TextureLoader().load(('/' + iconSensorReading))
        }
      },
      points: {
        position: new BufferAttribute(new Float32Array(0),3),
        color: new BufferAttribute(new Float32Array(0),4),
        index: new BufferAttribute(new Uint16Array(0),1)
      },
      line: {
        position: new BufferAttribute(new Float32Array(0),3),
        color: new BufferAttribute(new Float32Array(0),4),
        index: new BufferAttribute(new Uint16Array(0),1)
      }
    }
    for (let k in bufferGeometries) {
      for (let l in bufferGeometries[k]) {
        for (let i = 0; i < bufferGeometries[k][l].count; i++) {
          bufferGeometries[k][l].index.array[i] = i
        }
      }
    }

    let {width, height} = this.props

    this.state = {
      bufferGeometries,
      meshLine: new MeshLine(),
      meshLineMaterial: new MeshLineMaterial({lineWidth: 0.025, resolution: new Vector2(width,height), sizeAttenuation: 0, near: 1, far: 5000})
    }
  }

  componentWillReceiveProps (nextProps) {
    const { bufferGeometries, meshLine, meshLineMaterial } = this.state
    const {
      particles,
      pointsPath,
      readingPath
    } = nextProps
    const old_documents_length = this.props.pointsPath.length
    const new_documents_length = nextProps.pointsPath.length
    const new_reading_length = nextProps.readingPath.length

    if(old_documents_length !== new_documents_length){
        let indexes = new Uint16Array(new_reading_length)
        indexes = indexes.map((x,i) => i)

        bufferGeometries.line.position = new BufferAttribute(new Float32Array(new_reading_length * 3), 3)
        bufferGeometries.line.color = new BufferAttribute(new Float32Array(new_reading_length * 4), 4)
        bufferGeometries.line.index = new BufferAttribute(indexes,1)

        bufferGeometries.points.position = new BufferAttribute(new Float32Array(new_documents_length * 3), 3)
        bufferGeometries.points.color = new BufferAttribute(new Float32Array(new_documents_length * 4), 4)
        bufferGeometries.points.index = new BufferAttribute(new Uint16Array(new_documents_length),1)
    }

    let geometry =  new Geometry();
    console.log(1)
    readingPath.forEach((p) => {
        let [x,y,z,d] = p
        let v = new Vector3(x,y,15)
        geometry.vertices.push(v)
    })
     
    var line = new MeshLine();
    line.setGeometry(geometry)

    Object.keys(particles).forEach(type => {
      let uindexes = new Uint16Array(new_documents_length)
      let indexes = nextProps.pointsPath.map((x,i) => [x[2],i])
                                        .sort((b,a) => a[0] - b[0])
                                        .forEach((x,i) => uindexes[i] = x[1])

      bufferGeometries.line.position.array = particles[type].position
      bufferGeometries.line.position.needsUpdate = true
      bufferGeometries.line.color.array = particles[type].color
      bufferGeometries.line.color.needsUpdate = true

      bufferGeometries.points.position.array = particles[type].position
      bufferGeometries.points.position.needsUpdate = true
      bufferGeometries.points.color.array = particles[type].color
      bufferGeometries.points.color.needsUpdate = true
      bufferGeometries.points.index.array = uindexes
      bufferGeometries.points.index.needsUpdate = true
    })


    this.setState({
      bufferGeometries,
      meshLine: line
    })
  }

  render () {
    const { project } = ViewportMercator(this.props)
    const { 
      width,
      height,
      longitude,
      latitude
    } = this.props

    const { 
      bufferGeometries,
      meshLine,
      meshLineMaterial
     } = this.state

    const point = project([longitude, latitude])
    const startPoint = project([this.state.longitude, this.state.latitude])
    const left = point[0] - startPoint[0]
    const top = 0 - (point[1] - startPoint[1])
    const cameraProps = {
      left: 0,
      right: width,
      top: 0,
      bottom: height,
      near: 1,
      far: 5000,
      position: new Vector3(left, top, 600),
      lookAt: new Vector3(left, top, 0)
    }
    const readingPath = bufferGeometries.line
    const points = bufferGeometries.points
    const pixelRatio = window.devicePixelRatio ? window.devicePixelRatio : 1

    return (
      <div>
        <div id="three-renderer">
          <React3
            mainCamera="camera"
            width={ width }
            height={ height }
            alpha={ true }
            antialias={ true }
            pixelRatio={ pixelRatio }
          >
            <scene>
              <orthographicCamera
                name="camera"
                { ...cameraProps }
              />
              <mesh geometry={ meshLine.geometry } material={ meshLineMaterial }/>
              {
                Object.keys(bufferGeometries.particles).map(type => {
                  const particles = bufferGeometries.particles[type]
                  return (
                    <points key={ 'particles-' + type }>
                      <bufferGeometry
                        position={ points.position }
                        color={ points.color }
                        index={ points.index }
                      />
                      <pointsMaterial 
                        size={30}
                        vertexColors={VertexColors}
                        map={bufferGeometries.particles[type].texture}
                        transparent={true}
                      />
                    </points>
                  )
                })
              }
            </scene>
          </React3>
        </div>
      </div>
    )
  }
}

WebGLOverlay.propTypes = {}
