// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APINewTeam } from '../../api/types';
// import type { APIErrors, APINewExpedition } from '../../api/types';

type Props = {
  expeditionSlug: string,
  name?: string,
  slug?: string,
  description?: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (e: APINewTeam) => Promise<?APIErrors>;
}