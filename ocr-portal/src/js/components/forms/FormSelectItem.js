// @flow

import React, { Component } from 'react';

import { errorsFor } from '../../common/util';
import type { APIErrors } from '../../api/types';

type Props = {
labelText?: string;
name: string;
inline?: boolean;
className?: string;
ref?: string;
value: string | number | null;
firstOptionText: string;
options: {value: number | string, text: string}[];
errors: ?APIErrors;
onChange?: (e: any) => void;
}

export class FormSelectItem extends Component<void, $Exact<Props>, void> {
    render() {
        const {name, inline, className, labelText, ref, value, firstOptionText, options, errors, onChange} = this.props;

        return (
            <div className={ inline ? 'form-group inline' : 'form-group' }>
                { labelText &&
                  <label htmlFor={ name }>
                      { labelText }
                  </label> }
                <select name={ name } className={ className } value={ value } onChange={ onChange }>
                    <option value={ null } disabled>
                        { firstOptionText }
                    </option>
                    { options.map((o, i) => <option key={ i } value={ o.value }>
                                                { o.text }
                                            </option>) }
                </select>
                { errorsFor(errors, name) }
            </div>
            );
    }
}
