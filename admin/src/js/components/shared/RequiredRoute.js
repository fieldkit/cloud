// @flow weak

import React from 'react'
import { Route } from 'react-router-dom'
import ReactModal from 'react-modal'
import _ from 'lodash'

export const RouteOrLoading = ({component, required, path, ...rest}) => {
    const satisfied = !required || (_.isArrayLikeObject(required) && _.every(required));

    const renderFn = props => {
        if (satisfied) {
            return React.createElement(component, {
                ...props,
                ...rest
            });
        } else {
            return <div className="loading"></div>
        }
    }

    return <Route path={ path } render={ renderFn } />
};
