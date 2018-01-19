// @flow weak

import React, { Component } from 'react';
import _ from 'lodash';
import moment from 'moment';

const panelStyle: React.CSSProperties = {
    backgroundColor: '#f9f9f9',
    color: "#000",
    borderTopLeftRadius: 2,
    borderTopRightRadius: 2,
    borderBottomLeftRadius: 2,
    borderBottomRightRadius: 2,
    padding: '10px',
    // boxShadow: '0px 1px 4px rgba(0, 0, 0, .3)',
    // border: '1px solid rgba(0, 0, 0, 0.4)'
};

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    width: '300px',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
};

export default class FiltersPanel extends Component {
    props: {
        style: React.CSSProperties
    }

    render() {
        const { style, visibleFeatures } = this.props;
        const position = { top: 120, left: 70, bottom: 'auto', right: 'auto' };

        return (
            <div style={{ ...panelStyle, ...containerStyle, ...style, ...position }}>
                {this.renderSources(visibleFeatures.sources)}
            </div>
        );
    }

    renderSources(sources) {
        const { onShowFeature } = this.props;

        return (
            <div>
            <b>Sources</b>
            {_.values(sources).map(source => {
                const lastFeature = source.lastFeature;
                if (lastFeature == null) {
                    return (<div>Loading...</div>);
                }
                const lastFeatureDate = moment(new Date(lastFeature.properties.timestamp)).format('MMM Do YYYY, h:mm:ss a');
                const numberOfFeatures = source.numberOfFeatures;

                return (
                    <div key={source.id} style={{ paddingBottom: '5px' }}>
                        <div style={{ display: 'inline-block' }}>
                            <input type="checkbox" />
                        </div>

                        <div style={{ display: 'inline-block', paddingLeft: '10px' }}>
                            <a style={{color: 'black'}} href="#" onClick={() => console.log(source)}>{source.name}</a>
                            <span style={{ float: 'right' }}>&nbsp;{numberOfFeatures} features</span>
                        </div>
                        <div style={{ display: 'inline-block', paddingLeft: '30px' }}>
                            <a style={{color: 'black'}} href="#" onClick={() => onShowFeature(lastFeature) }>{lastFeatureDate}</a>
                        </div>
                    </div>
                )
            })}
            </div>
        );
    }
}
