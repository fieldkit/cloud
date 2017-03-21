
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import FKApiClient from '../api/api.js'

import I from 'immutable'

export const REQUEST_EXPEDITION = 'REQUEST_EXPEDITION'
export const INITIALIZE_EXPEDITION = 'INITIALIZE_EXPEDITION'
export const REQUEST_DOCUMENTS = 'REQUEST_DOCUMENTS'
export const INITIALIZE_DOCUMENTS = 'INITIALIZE_DOCUMENTS'
export const SET_VIEWPORT = 'SET_VIEWPORT'
export const UPDATE_DATE = 'UPDATE_DATE'
export const SELECT_PLAYBACK_MODE = 'SELECT_PLAYBACK_MODE'
export const SELECT_FOCUS_TYPE = 'SELECT_FOCUS_TYPE'
export const SELECT_ZOOM = 'SELECT_ZOOM'
export const JUMP_TO = 'JUMP_TO'
export const SET_MOUSE_POSITION = 'SET_MOUSE_POSITION'
export const SET_ZOOM = 'SET_ZOOM'

export function requestExpedition (expeditionID) {
  return function (dispatch, getState) {

    dispatch({
      type: REQUEST_EXPEDITION,
      id: expeditionID
    })

    let projectID = location.hostname.split('.')[0]
    if (projectID === 'localhost') projectID = 'eric'
    console.log('getting expedition')
    FKApiClient.getExpedition(projectID, expeditionID)
      .then(resExpedition => {
      // const resExpedition = {"name":"demoExpedition","slug":"demoexpedition"}
        console.log('expedition received:', resExpedition)
        if (!resExpedition) {
          console.log('error getting expedition')
        } else {
          console.log('expedition properly received')

          // {"name":"ian test","slug":"ian-test"}

          FKApiClient.getDocuments(projectID, expeditionID)
            .then(resDocuments => {
                // const resDocuments = [
                //   {"id":"GIUJ2WATF2HRCWBEKI34FUAVFZ4S7OHE","message_id":"PGMOCCWTIHR527UKISQ4EV4JYG55DJDU","request_id":"DZKYM26OJL2ENEQ5DG6VLDHMIXMTXQN2","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90714645385742,-77.04476165771484],"type":"Point"},"date":1.484660606e+09,"type":"Feature"}},{"id":"PI5KOJMYQPMA6VPUEPYGNUB2K3MULR45","message_id":"XGYLMLC5HHQE7LR3RG43LV5M4ER75C25","request_id":"53OTLO72IIXCTMN55PKUH3ZNCMTOD5KM","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90659713745117,-77.04178619384766],"type":"Point"},"date":1.484660868e+09,"type":"Feature"}},{"id":"UDTRVNZGAYTWEVJLJWV6NYG4SRFVGT4Z","message_id":"NLAJGOVGZQMCADGQEVD5JTXV3HAPBHR2","request_id":"SVYSPQGB37YPFNHI4U2COZEENV4O4ZEG","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.904422760009766,-77.04060363769531],"type":"Point"},"date":1.484661089e+09,"type":"Feature"}},{"id":"4HNYNXJVQYGCB6OEQVGH3AK6TKUQGZPS","message_id":"2TB3BCSU5Y6K2J3MTZ6NPZFFWTGLJOUK","request_id":"RXYNJ6IUP5FEAGXH5EAHUZMSOTPLPYI6","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90373611450195,-77.03819274902344],"type":"Point"},"date":1.484661273e+09,"type":"Feature"}},{"id":"W6CTJG6AGPNFLJATQRSISBWUIDE7GEQR","message_id":"QUL6V5QSLHJ2Y5CGKYUSC3NM3APHQCBD","request_id":"3KGKE44IVXRQOJKZERTSU5EHO44GTT4U","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484665769e+09,"geometry":{"coordinates":[38.903324127197266,-77.03483581542969],"type":"Point"},"type":"Feature"}},{"id":"4UPKWQKXF6DXS2WELGOKQGN2A6EBQRCP","message_id":"XM7IS7ADOSYVGOOGG2VCX5FOAJ2ZLF4R","request_id":"GHJR5SSJXWYZPZTQZQOUHQ5X2ZNHLHL2","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666355e+09,"geometry":{"coordinates":[38.90122985839844,-77.03242492675781],"type":"Point"},"type":"Feature"}},{"id":"IKQM4JVCYYRUUMKEJKU3337YGTYKCKDP","message_id":"XEY7FRYPKUOYMFGVI73Y7FMLPRNYG6MO","request_id":"HOJ6S3X7XU4HCNWPC3NFP5ZEMK57TMRO","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666632e+09,"geometry":{"coordinates":[38.898929595947266,-77.0322036743164],"type":"Point"},"type":"Feature"}},{"id":"IEO4AHCV5YG66ISODOXEMRWLJSQKRSSE","message_id":"ULJEQPA5EICPMLW3OBQDBE6FWVITEXC2","request_id":"NWIWIZWHQU4U3ZLLB2OPK6RDKQ6MDXF7","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666942e+09,"geometry":{"coordinates":[38.895423889160156,-77.03177642822266],"type":"Point"},"type":"Feature"}},{"id":"VWL2QQU4WHLVEXEJS44TYACS5NA2JDZC","message_id":"I4D4IKYJQ2H7PKM5U4QGFENZHJ7HEUZD","request_id":"O6T75DBTISNXCH6ZTWDBZ5EPYYCVUGUC","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666943e+09,"geometry":{"coordinates":[38.895423889160156,-77.03177642822266],"type":"Point"},"type":"Feature"}},{"id":"ZUEARF42RIEPXS4WF2SLIJPBDX3AVBSH","message_id":"76QTIEBEWQJREJ4KKKP65KSYIKNN5XMR","request_id":"FG6AVUVLLSAHS5JIONXJVVNXPZRNI7EU","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667222e+09,"geometry":{"coordinates":[38.8921012878418,-77.03201293945312],"type":"Point"},"type":"Feature"}},{"id":"F3I74F2IKMNYWC6UW5KXWHUD5BQRF4MR","message_id":"VRROPLG4GJB6QOSY4U6GKYH6RYCQDMF3","request_id":"R2UFNTJTV6GVY3S4T4ZZU24NFXRBUW7U","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667222e+09,"geometry":{"coordinates":[38.8921012878418,-77.03201293945312],"type":"Point"},"type":"Feature"}},{"id":"N7UF6ZMMUV7AZ47QLNUFCUSIJ56Q4472","message_id":"P5D553WE2XDZZHONZQHSEHP4GTQMIKV2","request_id":"U37WZSLHHW7WZP33TAX2V3XJGMSIIKWC","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667513e+09,"geometry":{"coordinates":[38.88996124267578,-77.03324127197266],"type":"Point"},"type":"Feature"}},{"id":"ANCFSUOIBB2X2O5HMDZD3EIBMFUVDSZV","message_id":"W2A6MND3WNAJIHK5IYUTVTFVDI3PH2FS","request_id":"QTORCPX2EXOLAHC5D4J3FJXTEYPO6QZY","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667654e+09,"geometry":{"coordinates":[38.89020919799805,-77.03443908691406],"type":"Point"},"type":"Feature"}},{"id":"JJJAL37D3LGNKRFTOSC7RYPGEE4SZVB7","message_id":"X3DJJMWFFCKBMING2H2PF5E2XAQYQBF4","request_id":"MAFNBRJTDYG4ZQQ7OGPGO5I5YWOO27KZ","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667654e+09,"geometry":{"coordinates":[38.89020919799805,-77.03443908691406],"type":"Point"},"type":"Feature"}},{"id":"CGSSZVB4FCFIBCKHE2MKTE2HS6SHCHS2","message_id":"LNJC7SAN2IEQC6Y6NBH7SBWT2T6LYN2H","request_id":"JKHV3NBO5FQZNH2JBHAUBRYGBI7WTFZ4","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667941e+09,"geometry":{"coordinates":[38.88999557495117,-77.036865234375],"type":"Point"},"type":"Feature"}},{"id":"E33FUBH4K43OI5YC3ZLACVAK7N3OV3OT","message_id":"KBQ2VWSHDMTZ2SNHDA7BUMG4TKYOPXJ5","request_id":"JW5IJEOHG622F6MIHKVTI3DHMIUOM6YU","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668077e+09,"geometry":{"coordinates":[38.889949798583984,-77.03812408447266],"type":"Point"},"type":"Feature"}},{"id":"XV34CCAQR4HS2QIMIZWDKXR755JL3TWR","message_id":"FRN7NUFLWUAXCTRPJEI5WSTQRSOWLANH","request_id":"TSI3UAFOXQJ27Q6PD3STXS7WHFMAMAA5","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668214e+09,"geometry":{"coordinates":[38.88996887207031,-77.03929901123047],"type":"Point"},"type":"Feature"}},{"id":"4RVD5JHXAJYEIO2UJNATMHG5ZZ443ZAG","message_id":"ASORVOLEAPKL6GQG6UPFMWTO5L7CLNOJ","request_id":"UJDATIJRFIENQNHO4J6FXAVAA44HKAXL","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668214e+09,"geometry":{"coordinates":[38.88996887207031,-77.03929901123047],"type":"Point"},"type":"Feature"}},{"id":"JSIDNPQVXZI76BUG2ZOK3CGORHHQFCWF","message_id":"NXXN7J5XB7TIN4YMQVICXCPJIJPXUMZY","request_id":"CJPM6TT4LFHFM3UOCRDZUSOQI3ADXDER","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668532e+09,"geometry":{"coordinates":[38.88978576660156,-77.04022979736328],"type":"Point"},"type":"Feature"}},{"id":"ZRWV3QVJTNM7EBDV74J7GEQFCGXEFZKS","message_id":"ELPFS6NVTUZYNC73X5X3CWN4BIFJU5SU","request_id":"AE5BX4Q4HK4WDBZFFNUFM5ZY2Q3CU2RR","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.48466868e+09,"geometry":{"coordinates":[38.889163970947266,-77.04084777832031],"type":"Point"},"type":"Feature"}},{"id":"PJUNXICZFRE43FEDY7QIH5QCAABKU2U5","message_id":"O6UXFWTZ3D6RV6S4LPPPT7OOCIW6CGZF","request_id":"G66JP4SOCJ2EN6FWWXW2YKV4DV6INDSQ","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668818e+09,"geometry":{"coordinates":[38.88915252685547,-77.04115295410156],"type":"Point"},"type":"Feature"}}]

                if (!resDocuments) resDocuments = []

              console.log('documents received:', resDocuments)
              if (!resDocuments) {
                console.log('error getting documents')
              } else {
                console.log('documents properly received')
                // const documents = resDocuments.map(d => {
                //   d.data = d.data * 1000
                //   return d
                //     .set('date', d.data.date * 1000)
                // })

                const documentMap = {}
                resDocuments
                  .forEach((d, i) => {
                    d.data.id = d.id
                    d.data.date = d.data.date * 1000
                    if (!d.data.geometry) d.data.geometry = d.data.Geometry
                    documentMap[d.id] = d.data
                  })
                const documents = I.fromJS(documentMap)

                if (documents.size > 0) {
                  const startDate = documents.toList()
                    .sort((d1, d2) => {
                      return d1.get('date') - d2.get('date')
                    })
                    .get(0).get('date')
                  const endDate = documents.toList()
                    .sort((d1, d2) => {
                      return d1.get('date') - d2.get('date')
                    })
                    .get(resDocuments.length - 1).get('date')

                  console.log(new Date(startDate), new Date(endDate))

                  const expeditionData = I.fromJS({
                    id: expeditionID,
                    name: expeditionID,
                    focusType: 'sensor-reading',
                    startDate,
                    endDate
                  })

                  dispatch([
                    {
                      type: INITIALIZE_EXPEDITION,
                      id: expeditionID,
                      data: expeditionData
                    },
                    {
                      type: INITIALIZE_DOCUMENTS,
                      data: documents
                    }
                  ])
                } else {
                  const startDate = new Date()
                  const endDate = new Date()
                  const expeditionData = I.fromJS({
                    id: expeditionID,
                    name: expeditionID,
                    focusType: 'sensor-reading',
                    startDate,
                    endDate
                  })

                  dispatch([
                    {
                      type: SET_ZOOM,
                      zoom: 2
                    },
                    {
                      type: INITIALIZE_EXPEDITION,
                      id: expeditionID,
                      data: expeditionData
                    },
                    {
                      type: INITIALIZE_DOCUMENTS,
                      data: documents
                    }
                  ])
                }
              }
            })
        }
      }) 
  }
}


export function requestExpeditions () {
  return function (dispatch, getState) {
    const projectID = location.hostname.split('.')[0]
    const expeditionID = getState().expeditions.get('currentExpeditionID')
    FKApiClient.getExpedition(projectID, expeditionID)
      .then(res => {
        console.log('expeditions received:', res)
        if (!res) {
          console.log('error getting expedition')
        } else {
          console.log('expedition properly received')
        }
      })
  }
}

export function initializeExpedition (id, data) {
  return function (dispatch, getState) {
    dispatch({
      type: INITIALIZE_EXPEDITION,
      id,
      data: I.fromJS({
        ...data,
        expeditionFetching: false,
        documentsFetching: true
      })
    })
  }
}

export function requestDocuments (id) {
  return function (dispatch, getState) {
    dispatch({
      type: REQUEST_DOCUMENTS
    })
    setTimeout(() => {
      const res = {
        'reading-0': {
          id: 'reading-0',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.6, 10.1]
          },
          date: 1484328718000
        },
        'reading-1': {
          id: 'reading-1',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.68, 10.11]
          },
          date: 1484328818000
        },
        'reading-2': {
          id: 'reading-2',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.7, 10.31]
          },
          date: 1484328958000
        },
        'reading-3': {
          id: 'reading-3',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.3, 10.35]
          },
          date: 1484329258000
        }
      }
      dispatch(initializeDocuments(id, I.fromJS(res)))
    }, 500)
  }
}

export function initializeDocuments (id, data) {
  return function (dispatch, getState) {
    dispatch({
      type: INITIALIZE_DOCUMENTS,
      id,
      data
    })
  }
}

export function setViewport(viewport, manual) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_VIEWPORT,
      viewport,
      manual
    })
  }
}

export function updateDate (date, playbackMode) {
  return function (dispatch, getState) {
    dispatch({
      type: UPDATE_DATE,
      date,
      playbackMode
    })
  }
}

export function selectPlaybackMode (mode) {
  return function (dispatch, getState) {
    dispatch({
      type: SELECT_PLAYBACK_MODE,
      mode
    })
  }
}

export function selectFocusType (focusType) {
  return function (dispatch, getState) {
    dispatch({
      type: SELECT_FOCUS_TYPE,
      focusType
    })
  }
}

export function selectZoom (zoom) {
  return function (dispatch, getState) {
    dispatch({
      type: SELECT_ZOOM,
      zoom
    })
  }
}

export function setMousePosition (x, y) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_MOUSE_POSITION,
      x,
      y
    })
  }
}
