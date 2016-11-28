
import * as actions from '../actions'
import * as d3 from 'd3'
import randomColor from 'randomcolor'

const okavangoReducer = (
  state = {
    mapStateNeedsUpdate: false,
    animate: false,
    isFetching: false,
    selectedExpedition: null,
    expeditions: {},
    speciesColors: {},
    contentActive: location.pathname.indexOf('/map') === -1,
    initialPage: location.pathname,
    lightBoxActive: false,
    lightBoxID: ''
  },
  action
) => {
  var expeditions, features, id, expeditionID, expedition, days

  switch (action.type) {
    case actions.CLOSE_LIGHTBOX:
      return {
        ...state,
        lightBoxActive: false,
        lightBoxID: ''
      }
    case actions.SHOW_360_PICTURE:
      return {
        ...state,
        lightBoxActive: true,
        lightBoxPost: action.post
      }
    case actions.HIDE_360_PICTURE:
      return {
        ...state,
        lightBoxActive: false,
        lightBoxID: ''
      }
    case actions.SET_PAGE:
      return {
        ...state,
        lightBoxActive: false
      }
    case actions.ENABLE_CONTENT:
      return {
        ...state,
        contentActive: true
      }
    case actions.RECEIVE_POSTS:
      // console.log('RECEIVED', action.data)
      expeditionID = action.expeditionID
      expedition = state.expeditions[expeditionID]

      // initializing days
      var timeRange = action.timeRange
      var postsByDay = {}
      var start = new Date(timeRange[0])
      var end = new Date(timeRange[1])
      var startDay = Math.floor((start.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
      var endDay = Math.floor((end.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
      for (i = startDay; i <= endDay; i++) {
        postsByDay[i] = {}
      }

      features = {}
      action.data.forEach((f) => {
        var id = f.id
        var flag = true
        if (!f.geometry) flag = false
        if (f.properties.FeatureType === 'image' && f.properties.Make === 'RICOH') flag = false
        if (flag) {
          features[id] = featureReducer(expedition.features[id], action, f)
        }
      })

      Object.keys(features).forEach((id) => {
        var feature = features[id]
        var day = Math.floor((new Date(feature.properties.DateTime).getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
        if (!postsByDay[day]) postsByDay[day] = {}
        postsByDay[day][id] = feature
      })
      Object.keys(postsByDay).forEach((k) => {
        postsByDay[k] = Object.assign({}, expedition.postsByDay[k], postsByDay[k])
      })

      return Object.assign({}, state, {
        // isFetchingPosts: state.isFetchingPosts - 1,
        mapStateNeedsUpdate: false,
        expeditions: Object.assign({}, state.expeditions, {
          [expeditionID]: Object.assign({}, expedition, {
            features: Object.assign({}, expedition.features, features),
            postsByDay: Object.assign({}, expedition.postsByDay, postsByDay)
          })
        })
      })

    case actions.FETCH_POSTS_BY_DAY:
      id = action.expeditionID
      return Object.assign({}, state, {
        // isFetchingPosts: state.isFetchingPosts + 1,
        expeditions: Object.assign({}, state.expeditions, {
          [id]: expeditionReducer(state.expeditions[id], action)
        })
      })

    case actions.RECEIVE_TOTAL_SIGHTINGS:
      id = action.id
      return Object.assign({}, state, {
        expeditions: Object.assign({}, state.expeditions, {
          [id]: Object.assign({}, state.expeditions[action.id], {
            totalSightings: action.data.results
          })
        })
      })

    case actions.COMPLETE_DAYS:
      id = action.expeditionID
      return Object.assign({}, state, {
        mapStateNeedsUpdate: false,
        expeditions: Object.assign({}, state.expeditions, {
          [id]: expeditionReducer(state.expeditions[id], action)
        })
      })

    case actions.SHOW_LOADING_WHEEL:
      return Object.assign({}, state, {
        mapStateNeedsUpdate: false,
        isFetching: true
      })

    case actions.HIDE_LOADING_WHEEL:
      return Object.assign({}, state, {
        mapStateNeedsUpdate: true,
        isFetching: false
      })

    case actions.START:
      return Object.assign({}, state, {
        mapStateNeedsUpdate: true,
        animate: true
      })

    case actions.UPDATE_TIME:
      expeditionID = action.expeditionID
      return Object.assign({}, state, {
        mapStateNeedsUpdate: action.updateMapState,
        expeditions: Object.assign({}, state.expeditions, {
          [expeditionID]: expeditionReducer(state.expeditions[expeditionID], action)
        })
      })

    case actions.UPDATE_MAP:
      expeditionID = action.expeditionID
      return Object.assign({}, state, {
        mapStateNeedsUpdate: true,
        expeditions: Object.assign({}, state.expeditions, {
          [expeditionID]: expeditionReducer(state.expeditions[expeditionID], action)
        })
      })

    case actions.RECEIVE_EXPEDITIONS:
      expeditions = {}
      var latestDate = new Date(0)
      var latestExpedition
      // console.log('aga', location.pathname.split('/'))
      Object.keys(action.data.results).forEach((id) => {
        var e = action.data.results[id]
        expeditions[id] = expeditionReducer(state.expeditions[id], action, e)
        if (expeditions[id].start.getTime() + expeditions[id].dayCount * (1000 * 3600 * 24) > latestDate.getTime()) {
          latestDate = new Date(expeditions[id].start.getTime() + expeditions[id].dayCount * (1000 * 3600 * 24))
          latestExpedition = id
        }
      })

      return Object.assign({}, state, {
        expeditions: Object.assign({}, state.expeditions, expeditions),
        isFetching: false,
        selectedExpedition: latestExpedition
      })

    case actions.SET_EXPEDITION:
      var selectedExpedition = action.id
      return Object.assign({}, state, {
        mapStateNeedsUpdate: true,
        selectedExpedition: selectedExpedition
      })

    case actions.SET_CONTROL:
      id = state.selectedExpedition
      expeditions = {
        [id]: expeditionReducer(state.expeditions[id], action)
      }
      return Object.assign({}, state, {
        mapStateNeedsUpdate: action.target === 'zoom',
        expeditions: Object.assign({}, state.expeditions, expeditions)
      })

    case actions.RECEIVE_DAY:
      expeditionID = action.expeditionID
      expedition = state.expeditions[expeditionID]

      // initialize feature buckets
      var members = { ...expedition.members }
      var featuresByMember = {}
      var featuresByDay = {}
      var ambitsByTile = {}
      var dateRange = action.dateRange.map((d) => {
        return Math.floor((new Date(d).getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
      })
      for (var i = dateRange[0]; i <= dateRange[1]; i++) {
        featuresByDay[i] = {}
      }

      // initialize features
      features = {}
      action.data.forEach((f) => {
        var id = f.id
        if (f.properties.Team === 'RiverMain') {
          if(f.properties.FeatureType !== 'ambit_geo' || f.properties.Member === 'Steve' || f.properties.Member === 'GB' || f.properties.Member === 'Jer' || f.properties.Member === 'Shah'){
            features[id] = featureReducer(expedition.features[id], action, f)
            if (f.properties.FeatureType === 'ambit_geo') {
              if (!members[f.properties.Member]) {
                members[f.properties.Member] = {
                  color: expedition.memberColors[d3.values(members).length % expedition.memberColors.length]
                }
              }
            }
          }
        }
      })

      var tileResolution = Math.floor((expedition.geoBounds[2] - expedition.geoBounds[0]) * 111 / 10)
      var coordinatesToTile = (coordinates, geoBounds) => {
        var x = Math.floor((coordinates[0] - geoBounds[0]) * 111 / 10)
        var y = Math.floor((coordinates[1] - geoBounds[3]) * 111 / 10)
        return {x, y}
      }

      // assign feature to day, tile and member
      Object.keys(features).forEach((id) => {
        var feature = features[id]
        // assign feature to day
        var day = Math.floor((new Date(feature.properties.DateTime).getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
        var type = feature.properties.FeatureType
        if (!featuresByDay[day]) featuresByDay[day] = {}
        if (!featuresByDay[day][type]) featuresByDay[day][type] = {}
        featuresByDay[day][type][id] = feature

        if (feature.properties.FeatureType === 'ambit_geo') {
          // assign feature to member
          var memberID = feature.properties.Member
          if (!members[memberID]) {
            members[memberID] = {
              name: memberID,
              color: expedition.memberColors[d3.values(members).length % expedition.memberColors.length]
            }
          }
          var dayID = Math.floor((new Date(feature.properties.DateTime).getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
          if (!featuresByMember[memberID]) featuresByMember[memberID] = {}
          if (!featuresByMember[memberID][dayID]) featuresByMember[memberID][dayID] = {}
          featuresByMember[memberID][dayID][id] = feature

          // assign feature to tile
          var tileCoordinates = coordinatesToTile(feature.geometry.coordinates, expedition.geoBounds)
          var tileID = tileCoordinates.x + tileCoordinates.y * tileResolution
          if (!ambitsByTile[tileID]) ambitsByTile[tileID] = {}
          ambitsByTile[tileID][id] = feature
        }
      })

      const extendFeatures = (bucket) => {
        if (d3.values(bucket).length > 0) {
          // pick the two earliest and latest features
          var timeRange = [new Date(), new Date(0)]
          var featureRange = []
          d3.values(bucket).forEach((f) => {
            var dateTime = new Date(f.properties.DateTime)
            if (timeRange[0].getTime() > dateTime.getTime()) {
              timeRange[0] = dateTime
              featureRange[0] = f
            }
            if (timeRange[1].getTime() < dateTime.getTime()) {
              timeRange[1] = dateTime
              featureRange[1] = f
            }
          })

          // clone features with new dates
          var start = new Date(timeRange[0].getTime() - (timeRange[0].getTime() % (1000 * 3600 * 24)))
          var end = new Date(start.getTime() + (1000 * 3600 * 24))
          id = Date.now() + (Math.floor(Math.random() * 10000) / 10000)
          bucket[id] = Object.assign({}, featureRange[0])
          bucket[id].properties = Object.assign({}, bucket[id].properties, {
            DateTime: start.toString()
          })
          id = Date.now() + (Math.floor(Math.random() * 10000) / 10000)
          bucket[id] = Object.assign({}, featureRange[1])
          bucket[id].properties = Object.assign({}, bucket[id].properties, {
            DateTime: end.toString()
          })
        }
      }
      Object.keys(featuresByDay).forEach(d => {
        featuresByDay[d] = Object.assign({}, expedition.featuresByDay[d], featuresByDay[d])
        extendFeatures(featuresByDay[d].beacon)
      })
      Object.keys(featuresByMember).forEach(m => {
        Object.keys(featuresByMember[m]).forEach(d => {
          extendFeatures(featuresByMember[m][d])
        })
      })

      featuresByDay = Object.assign({}, expedition.featuresByDay, featuresByDay)
      days = Object.assign({}, featuresByDay)
      for (var d in days) {
        days[d] = dayReducer(expedition.days[d], action, featuresByDay[d])
      }

      Object.keys(featuresByMember).forEach((k) => {
        featuresByMember[k] = Object.assign({}, expedition.featuresByMember[k], featuresByMember[k])
      })

      Object.keys(ambitsByTile).forEach((k) => {
        ambitsByTile[k] = Object.assign({}, expedition.ambitsByTile[k], ambitsByTile[k])
      })

      return Object.assign({}, state, {
        mapStateNeedsUpdate: false,
        expeditions: Object.assign({}, state.expeditions, {
          [expeditionID]: Object.assign({}, expedition, {
            days: Object.assign({}, expedition.days, days),
            // features: Object.assign({}, expedition.features, features),
            featuresByDay: featuresByDay,
            featuresByMember: Object.assign({}, expedition.featuresByMember, featuresByMember),
            ambitsByTile: Object.assign({}, expedition.ambitsByTile, ambitsByTile),
            members
          })
        })
      })

    case actions.RECEIVE_FEATURES:
      expeditionID = action.expeditionID
      expedition = state.expeditions[expeditionID]

      var tileRange = action.tileRange
      var featuresByTile = {}
      tileRange.forEach((t, i) => {
        featuresByTile[t] = {}
      })

      features = {}
      action.data.forEach((f) => {
        var id = f.id
        if (f.properties.Team === 'RiverMain') {
          var flag = true
          if (f.properties.FeatureType === 'sighting') {
            if (!f.properties.Taxonomy) f.properties.color = 0xb4b4b4
            else {
              var taxClass = f.properties.Taxonomy.Class
              if (!state.speciesColors[taxClass]) {
                state.speciesColors[taxClass] = parseInt(randomColor({ luminosity: 'light', format: 'hex' }).slice(1), 16)
              }
              f.properties.color = state.speciesColors[taxClass]
            }
          }
          if (f.properties.FeatureType === 'tweet' && f.properties.Text && f.properties.Text[0] === '@') flag = false
          // if (f.properties.FeatureType === 'image' && f.properties.Make === 'RICOH') flag = false
          if (flag) {
            features[id] = featureReducer(expedition.features[id], action, f)
          }
        }
      })

      var tileResolution = Math.floor((expedition.geoBounds[2] - expedition.geoBounds[0]) * 111 / 10)
      var coordinatesToTile = (coordinates, geoBounds) => {
        var x = Math.floor((coordinates[0] - geoBounds[0]) * 111 / 10)
        var y = Math.floor((coordinates[1] - geoBounds[3]) * 111 / 10)
        return {x, y}
      }

      Object.keys(features).forEach((id) => {
        var feature = features[id]
        var tileCoordinates = coordinatesToTile(feature.geometry.coordinates, expedition.geoBounds)
        var tileID = tileCoordinates.x + tileCoordinates.y * tileResolution
        if (!featuresByTile[tileID]) featuresByTile[tileID] = {}
        featuresByTile[tileID][id] = feature
      })
      Object.keys(featuresByTile).forEach((k) => {
        featuresByTile[k] = Object.assign({}, expedition.featuresByTile[k], featuresByTile[k])
      })

      for (var k in features) {
        var feature = features[k]
        if (feature.properties.FeatureType === 'sighting') {
          delete features[k]
        }
      }

      return Object.assign({}, state, {
        mapStateNeedsUpdate: false,
        expeditions: Object.assign({}, state.expeditions, {
          [expeditionID]: Object.assign({}, expedition, {
            features: Object.assign({}, expedition.features, features),
            featuresByTile: Object.assign({}, expedition.featuresByTile, featuresByTile)
          })
        })
      })

    case actions.SELECT_FEATURE:
      break

    case actions.UNSELECT_FEATURE:
      break
    default:
      return state
  }

  return state
}

const expeditionReducer = (
  state = {
    name: '',
    playback: 'forward',
    layout: 'rows',
    initialZoom: 4,
    targetZoom: 15,
    zoom: 15,
    isFetching: false,
    geoBounds: [-8, -21.5, 25.5, 12],
    tileSize: 10,
    start: new Date(),
    end: new Date(0),
    currentDate: new Date(),
    dayCount: 0,
    days: [],
    features: {},
    featuresByTile: {},
    ambitsByTile: {},
    featuresByDay: {},
    postsByDay: {},
    featuresByMember: {},
    mainFocus: 'Explorers',
    secondaryFocus: 'Steve',
    coordinates: [0, 0],
    current360Images: [],
    currentGeoBounds: [0, 0],
    currentPosts: [],
    currentSightings: [],
    currentAmbits: [],
    totalSightings: [],
    members: {},
    currentMembers: [],
    memberColors: [
      0xFDBF6F,
      0xA6CEE3,
      0xB2DF8A,
      0xFB9A99,
      0xCAB2D6,
      0xFCEA97,
      0xB4F0D1,
      0xBFBFFF,
      0xFFABD5
    ]
    // memberColors: [
    //   'rgba(253, 191, 111, 1)',
    //   'rgba(166, 206, 227, 1)',
    //   'rgba(178, 223, 138, 1)',
    //   'rgba(251, 154, 153, 1)',
    //   'rgba(202, 178, 214, 1)',
    //   'rgba(252, 234, 151, 1)',
    //   'rgba(180, 240, 209, 1)',
    //   'rgba(191, 191, 255, 1)',
    //   'rgba(255, 171, 213, 1)'
    // ],
  },
  action,
  data
) => {
  var i
  switch (action.type) {
    case actions.FETCH_POSTS_BY_DAY:
      var postsByDay = []
      // var start = new Date(action.range[0])
      // var end = new Date(action.range[1])
      // var startDay = Math.floor((start.getTime() - state.start.getTime()) / (1000 * 3600 * 24))
      // var endDay = Math.floor((end.getTime() - state.start.getTime()) / (1000 * 3600 * 24))
      // for (i = startDay; i <= endDay; i++) {
      action.daysToFetch.forEach((d) => {
        postsByDay[d] = 'loading'
      })
      return Object.assign({}, state, {
        postsByDay: Object.assign({}, state.postsByDay, postsByDay)
      })

    case actions.COMPLETE_DAYS:
      var days = Object.assign({}, state.days)

      // add mock days at both ends of the expedition
      for (i = 0; i < state.dayCount; i++) {
        if (days[i] && !days[i].incomplete) {
          days[-1] = Object.assign({}, days[i])
          break
        }
      }
      for (i = state.dayCount; i >= 0; i--) {
        if (days[i] && !days[i].incomplete) {
          days[state.dayCount] = Object.assign({}, days[i])
          break
        }
      }

      // detect incomplete days
      var incompleteRange = [-1, -1]
      var completedDays = []
      for (i = 0; i < state.dayCount; i++) {
        var day = days[i]
        if (!day) {
          incompleteRange = [-1, -1]
        } else {
          if (day.incomplete && days[i - 1] && !days[i - 1].incomplete) {
            incompleteRange[0] = i
            incompleteRange[1] = -1
          }
          if (day.incomplete && days[i + 1] && !days[i + 1].incomplete) {
            incompleteRange[1] = i
          }
        }

        // full data gap detected, filling in
        if (incompleteRange[0] > -1 && incompleteRange[1] > -1) {
          // look for filling values bordering the gap
          var fillingDays = [days[+incompleteRange[0] - 1], days[+incompleteRange[1] + 1]]
          var fillingBeacons = [
            d3.values(fillingDays[0].beacons).slice(0).sort((a, b) => {
              return new Date(b.properties.DateTime).getTime() - new Date(a.properties.DateTime).getTime()
            })[0],
            d3.values(fillingDays[1].beacons).slice(0).sort((a, b) => {
              return new Date(a.properties.DateTime).getTime() - new Date(b.properties.DateTime).getTime()
            })[0]
          ]

          // fill in gaps
          var l2 = Math.ceil((incompleteRange[1] - incompleteRange[0] + 1) / 2)
          for (var j = 0; j < l2; j++) {
            var dayIndex = [(+incompleteRange[0] + j), (+incompleteRange[1] - j)]
            for (var k = 0; k < 2; k++) {
              for (var l = 0; l < 2; l++) {
                // here k === 0 removes gradual translation between day 1 to day 2
                if (k === 0 || (days[dayIndex[0]] !== days[dayIndex[1]]) || l === 0) {
                  var dayID = dayIndex[k]
                  day = days[dayID]
                  var date = new Date(state.start.getTime() + (1000 * 3600 * 24) * (dayID + (k === l ? 0 : 1)))
                  var beaconID = Date.now() + (Math.floor(Math.random() * 10000) / 10000)
                  day.beacons[beaconID] = Object.assign({}, fillingBeacons[k])
                  day.beacons[beaconID].properties = Object.assign({}, day.beacons[beaconID].properties, {
                    DateTime: date
                  })
                  day.incomplete = false
                  if (completedDays.indexOf(dayID) === -1) completedDays.push(dayID)
                }
              }
            }
          }
          incompleteRange = [-1, -1]
        }
      }

      // remove mock days at both ends of the expedition
      delete days[-1]
      delete days[state.dayCount]

      // console.log('fill following days:', completedDays, days)
      return Object.assign({}, state, {
        days: days
      })

    case actions.UPDATE_TIME:
      return Object.assign({}, state, {
        currentDate: action.currentDate
      })

    case actions.UPDATE_MAP:

      // initializing featuresByTile entries so they won't be queried multiple times
      action.tileRange.forEach(t => {
        if (!state.featuresByTile[t]) state.featuresByTile[t] = {}
      })

      var currentSightings = []
      var current360Images = []
      var currentPosts = []
      var currentDays = []
      var currentAmbits = {}
      var currentMembers = []

      var allAmbits = []
      action.tilesInView.forEach((t) => {
        // sort features by type
        var features = {}
        d3.values(state.featuresByTile[t]).forEach((f) => {
          if (!features[f.properties.FeatureType]) features[f.properties.FeatureType] = []
          features[f.properties.FeatureType].push(f)
          var day = Math.floor((new Date(f.properties.DateTime).getTime() - state.start.getTime()) / (1000 * 3600 * 24))
          if (currentDays.indexOf(day) === -1) currentDays.push(day)
        })

        // this def could be written more elegantly...
        d3.values(state.ambitsByTile[t]).forEach((f) => {
          var day = Math.floor((new Date(f.properties.DateTime).getTime() - state.start.getTime()) / (1000 * 3600 * 24))
          if (currentDays.indexOf(day) === -1) currentDays.push(day)
        })

        if (features.sighting) {
          var sightings = features.sighting.map((f) => {
            return {
              position: {
                x: f.geometry.coordinates[0],
                y: f.geometry.coordinates[1],
                z: 0
              },
              radius: f.properties.radius,
              color: f.properties.color,
              type: f.properties.FeatureType,
              date: new Date(f.properties.DateTime),
              name: f.properties.SpeciesName,
              count: f.properties.Count
            }
          })
          currentSightings = currentSightings.concat(sightings)
        }

        if (features.image) {
          var images = features.image.filter(i => {
            return i.properties.Make === 'RICOH'
          })
          current360Images = current360Images.concat(images)
        }

        current360Images.forEach((image, i) => {
          image.properties.next = current360Images[i - 1]
          image.properties.previous = current360Images[i + 1]
        })

        var allPosts = []
        if (features.tweet) allPosts = allPosts.concat(features.tweet)
        if (features.audio) allPosts = allPosts.concat(features.audio)
        if (features.blog) allPosts = allPosts.concat(features.blog)
        if (features.image) allPosts = allPosts.concat(features.image)
        if (allPosts) {
          var posts = allPosts.map((f) => {
            return {
              position: [
                f.geometry.coordinates[0],
                f.geometry.coordinates[1]
              ],
              type: f.properties.FeatureType,
              id: f.id,
              properties: f.properties
            }
          })
          currentPosts = currentPosts.concat(posts)
        }
      })

      var paddingDays = []
      currentDays.forEach(d => {
        var flag = false
        for (i = d - 1; i <= d + 1; i++) {
          if (currentDays.indexOf(i) === -1 && paddingDays.indexOf(i) === -1) {
            paddingDays.push(i)
            flag = true
          }
        }
        if (flag) actions.fetchDay(new Date(state.start.getTime() + (1000 * 3600 * 24) * d))
      })
      currentDays = currentDays
        .concat(paddingDays)

      currentDays = currentDays
        .sort((a, b) => {
          return a - b
        })
        .forEach(d => {
          if (state.featuresByDay[d]) {
            allAmbits = allAmbits.concat(d3.values(state.featuresByDay[d].ambit_geo))
          }
        })

      allAmbits
        .sort((a, b) => {
          return new Date(a.properties.DateTime).getTime() - new Date(b.properties.DateTime).getTime()
        })
        .forEach(f => {
          var name = f.properties.Member
          if (!currentAmbits[name]) {
            currentAmbits[name] = {
              color: state.members[name].color,
              coordinates: [],
              dates: []
            }
          }
          if (!currentMembers[name]) currentMembers[name] = {}
          currentAmbits[name].coordinates.push(f.geometry.coordinates)
          currentAmbits[name].dates.push(f.properties.DateTime)
        })

      currentAmbits = d3.values(currentAmbits)

      return Object.assign({}, state, {
        coordinates: action.coordinates,
        currentAmbits,
        currentDate: action.currentDate,
        currentGeoBounds: action.viewGeoBounds,
        currentMembers,
        currentPosts,
        currentSightings,
        current360Images,
        zoom: action.zoom
      })

    case actions.RECEIVE_EXPEDITIONS:
      var dayCount = data.Days + 1
      // removing +1 here because we receive beacons before any other features on current day
      // var dayCount = data.Days
      var start = new Date(new Date(data.StartDate).getTime() + 2 * (1000 * 3600))
      var end = new Date(start.getTime() + dayCount * (1000 * 3600 * 24))
      // currentDate is 2 days before last beacon
      var currentDate = new Date(end.getTime() - 2 * (1000 * 3600 * 24))

      var name = data.Name

      // 111 km per latitude degree
      // ~ 10km per screen at zoom level 14
      // [west, north, east, south]
      var geoBounds = data.GeoBounds
      // var geoBoundsDistance = [(geoBounds[2] - geoBounds[0]) * 111, (geoBounds[3] - geoBounds[1]) * 111]

      return Object.assign({}, state, {
        name: name,
        start: start,
        currentDate: currentDate,
        end: end,
        dayCount: dayCount,
        geoBounds: geoBounds
      })

    case actions.SET_CONTROL:
      if (action.target === 'zoom') {
        if (action.mode === 'increment') action.mode = Math.max(1, Math.min(15, state.zoom + 1))
        if (action.mode === 'decrement') action.mode = Math.max(1, Math.min(15, state.zoom - 1))
      }
      return Object.assign({}, state, {
        [action.target]: action.mode
      })

    default:
      break
  }
  return state
}

const dayReducer = (
  state = {
    isFetching: false,
    start: new Date(),
    end: new Date(0),
    beacons: {},
    ambits: {},
    incomplete: true
  },
  action,
  features
) => {
  var start, end
  switch (action.type) {
    case actions.RECEIVE_DAY:
      start = new Date()
      end = new Date(0)
      if (!features.beacon) break
      var incomplete = Object.keys(features.beacon).length === 0

      Object.keys(features.beacon).forEach((k) => {
        var f = features.beacon[k]
        var d = new Date(f.properties.DateTime)
        if (d.getTime() < start.getTime()) start = d
        if (d.getTime() > end.getTime()) end = d
      })

      return Object.assign({}, state, {
        isFetching: false,
        start: start,
        end: end,
        beacons: Object.assign({}, state.beacons, features.beacon),
        ambits: Object.assign({}, state.ambits, features.ambit),
        incomplete: incomplete
      })

    default:
      break
  }

  return state
}

const featureReducer = (
  state = {},
  action,
  feature
) => {
  switch (action.type) {
    case actions.RECEIVE_POSTS:
      feature.properties.scatter = [((Math.random() * 2) - 1) * 0.00075, ((Math.random() * 2) - 1) * 0.00075]
      return Object.assign({}, state, feature)
    case actions.RECEIVE_DAY:
      return Object.assign({}, state, feature)
    case actions.RECEIVE_FEATURES:
      feature.properties.scatter = [((Math.random() * 2) - 1) * 0.00075, ((Math.random() * 2) - 1) * 0.00075]
      if (feature.properties.FeatureType === 'sighting') {
        feature.properties.radius = 2 + Math.sqrt(feature.properties.Count) * 2
      }
      return Object.assign({}, state, feature)
    default:
      break
  }

  return state
}

export default okavangoReducer
