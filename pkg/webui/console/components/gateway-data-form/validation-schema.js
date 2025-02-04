// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import Yup from '@ttn-lw/lib/yup'
import sharedMessages from '@ttn-lw/lib/shared-messages'
import { id as gatewayIdRegexp } from '@ttn-lw/lib/regexp'

import {
  attributeValidCheck,
  attributeTooShortCheck,
  attributeKeyTooLongCheck,
  attributeValueTooLongCheck,
} from '@console/lib/attributes'
import {
  addressWithOptionalScheme as addressWithOptionalSchemeRegexp,
  delay as delayRegexp,
} from '@console/lib/regexp'

const validationSchema = Yup.object().shape({
  owner_id: Yup.string(),
  ids: Yup.object().shape({
    gateway_id: Yup.string()
      .min(3, Yup.passValues(sharedMessages.validateTooShort))
      .max(36, Yup.passValues(sharedMessages.validateTooLong))
      .matches(gatewayIdRegexp, Yup.passValues(sharedMessages.validateIdFormat))
      .required(sharedMessages.validateRequired),
    eui: Yup.nullableString().length(8 * 2, Yup.passValues(sharedMessages.validateLength)),
  }),
  name: Yup.string()
    .min(2, Yup.passValues(sharedMessages.validateTooShort))
    .max(50, Yup.passValues(sharedMessages.validateTooLong)),
  update_channel: Yup.string()
    .min(2, Yup.passValues(sharedMessages.validateTooShort))
    .max(128, Yup.passValues(sharedMessages.validateTooLong)),
  description: Yup.string().max(2000, Yup.passValues(sharedMessages.validateTooLong)),
  frequency_plan_id: Yup.string().when(['$gsEnabled'], (gsEnabled, schema) => {
    if (!gsEnabled) {
      return schema.strip()
    }

    return schema
      .max(64, Yup.passValues(sharedMessages.validateTooLong))
      .required(sharedMessages.validateRequired)
  }),
  gateway_server_address: Yup.string().matches(
    addressWithOptionalSchemeRegexp,
    Yup.passValues(sharedMessages.validateAddressFormat),
  ),
  require_authenticated_connection: Yup.boolean().default(false),
  location_public: Yup.boolean().default(false),
  status_public: Yup.boolean().default(false),
  schedule_downlink_late: Yup.boolean().default(false),
  update_location_from_status: Yup.boolean().default(false),
  auto_update: Yup.boolean().default(false),
  schedule_anytime_delay: Yup.string().matches(
    delayRegexp,
    Yup.passValues(sharedMessages.validateDelayFormat),
  ),
  attributes: Yup.array()
    .max(10, Yup.passValues(sharedMessages.attributesValidateTooMany))
    .test(
      'has no empty string values',
      sharedMessages.attributesValidateRequired,
      attributeValidCheck,
    )
    .test(
      'has key length longer than 2',
      sharedMessages.attributeKeyValidateTooShort,
      attributeTooShortCheck,
    )
    .test(
      'has key length less than 36',
      sharedMessages.attributeKeyValidateTooLong,
      attributeKeyTooLongCheck,
    )
    .test(
      'has value length less than 200',
      sharedMessages.attributeValueValidateTooLong,
      attributeValueTooLongCheck,
    ),
})

export default validationSchema
