import React, {PropTypes, Component} from 'react'
import ViewportMercator from 'viewport-mercator-project'
import THREE from '../../vendors/react-three-renderer/node_modules/three'
import React3 from '../../vendors/react-three-renderer'
import sightingTexturePath from '../../img/sighting.png'

export default class WebGLOverlay extends Component {
  constructor (props) {
    super(props)

    const shaders = {
      particles: {
        vertexShader: [
          'attribute vec4 color;',
          'varying vec4 vColor;',
          'void main() {',
          '    vColor = color;',
          '    vec4 mvPosition = modelViewMatrix * vec4( position.xy, 0.0 , 1.0 );',
          '    gl_PointSize = float( position.z );',
          '    gl_Position = projectionMatrix * mvPosition;',
          '}'
        ].join('\n'),
        fragmentShader: [
          'varying vec4 vColor;',
          'uniform sampler2D texture;',
          'vec4 vFragColor;',
          'void main() {',
          '    vFragColor = vColor * texture2D( texture, gl_PointCoord );',
          '    if (vFragColor.w > 0.25) {',
          '      gl_FragColor = vFragColor;',
          '    } else {',
          '      discard;',
          '    }',
          '}'
        ].join('\n')
      }
    }

    const bufferGeometries = {
      focusParticles: {
        count: 1,
        position: new THREE.BufferAttribute(new Float32Array(1 * 3), 3),
        color: new THREE.BufferAttribute(new Float32Array(1 * 4), 4),
        index: new THREE.BufferAttribute(new Uint16Array(1 * 1), 1)
      },
      readingParticles: {
        count: 1000,
        position: new THREE.BufferAttribute(new Float32Array(1000 * 3), 3),
        color: new THREE.BufferAttribute(new Float32Array(1000 * 4), 4),
        index: new THREE.BufferAttribute(new Uint16Array(1000 * 1), 1)
      }
    }
    for (let k in bufferGeometries) {
      for (let i = 0; i < bufferGeometries[k].count; i++) {
        bufferGeometries[k].index.array[i] = i
      }
    }

    const sightingTexture = new THREE.TextureLoader().load(sightingTexturePath)

    this.state = {
      shaders,
      sightingTexture,
      bufferGeometries
    }
  }

  componentWillReceiveProps (nextProps) {
    const { bufferGeometries } = this.state
    const {
      focusParticles,
      readingParticles,
    } = nextProps

    bufferGeometries.focusParticles.position.array = focusParticles.position
    bufferGeometries.focusParticles.position.needsUpdate = true
    bufferGeometries.focusParticles.color.array = focusParticles.color
    bufferGeometries.focusParticles.color.needsUpdate = true
    
    bufferGeometries.readingParticles.position.array = readingParticles.position
    bufferGeometries.readingParticles.position.needsUpdate = true
    bufferGeometries.readingParticles.color.array = readingParticles.color
    bufferGeometries.readingParticles.color.needsUpdate = true

    this.setState({
      ...this.state,
      bufferGeometries
    })
  }

  render () {
    const { project } = ViewportMercator(this.props)

    const { 
      width,
      height,
      longitude,
      latitude,
      readingPath
    } = this.props

    const { 
      bufferGeometries,
      shaders,
      sightingTexture
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
      position: new THREE.Vector3(left, top, 600),
      lookAt: new THREE.Vector3(left, top, 0)
    }

    return (
      <div>
        <div id="three-renderer">
          <React3
            mainCamera="camera"
            width={width}
            height={height}
            onAnimate={this._onAnimate}
            alpha={true}
            antialias={true}
          >
            <scene>
              <orthographicCamera
                name="camera"
                { ...cameraProps }
              />
              {
                !!readingPath &&
                <line>
                  <geometry
                    vertices={readingPath}
                    dynamic={true}
                  >
                  </geometry>
                  <lineBasicMaterial
                    linewidth={10}
                    opacity={0.7}
                    transparent={false}
                    color={new THREE.Color('#ffffff')}
                  >
                  </lineBasicMaterial>
                </line>
              }
              { 
                !!bufferGeometries.readingParticles &&
                <points>
                  <bufferGeometry
                    position={ bufferGeometries.readingParticles.position }
                    color={ bufferGeometries.readingParticles.color }
                    index={ bufferGeometries.readingParticles.index }
                  />
                  <shaderMaterial
                    vertexShader={ shaders.particles.vertexShader}
                    fragmentShader={ shaders.particles.fragmentShader}
                    uniforms={{
                      texture: { 
                        type: 't',
                        value: sightingTexture 
                      }
                    }}
                  >
                  </shaderMaterial>
                </points>
              }
              { 
                !!bufferGeometries.focusParticles &&
                <points>
                  <bufferGeometry
                    position={ bufferGeometries.focusParticles.position }
                    color={ bufferGeometries.focusParticles.color }
                    index={ bufferGeometries.focusParticles.index }
                  />
                  <shaderMaterial
                    vertexShader={ shaders.particles.vertexShader}
                    fragmentShader={ shaders.particles.fragmentShader}
                    uniforms={{
                      texture: { 
                        type: 't',
                        value: sightingTexture 
                      }
                    }}
                  >
                  </shaderMaterial>
                </points>
              }
            </scene>
          </React3>
        </div>
      </div>
    )
  }
}

WebGLOverlay.propTypes = {}
