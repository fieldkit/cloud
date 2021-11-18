import _ from "lodash";
import moment from "moment";
import Vue from "vue";
import VueI18n, { LocaleMessages } from "vue-i18n";
import Messages from "@messageformat/runtime/messages";
import MessageFormat from "@messageformat/core";

class MessageFormatFormatter {
    //
    // interpolate
    //
    // @param {string} message
    //   string of list or named format.
    //   e.g.
    //   - named formatting: 'Hi {name}'
    //   - list formatting: 'Hi {0}'
    //
    // @param {Object | Array} values
    //   values of `message` interpolation.
    //   passed values with `$t`, `$tc` and `i18n` functional component.
    //   e.g.
    //   - $t('hello', { name: 'kazupon' }) -> passed values: Object `{ name: 'kazupon' }`
    //   - $t('hello', ['kazupon']) -> passed values: Array `['kazupon']`
    //   - `i18n` functional component (component interpolation)
    //     <i18n path="hello">
    //       <p>kazupon</p>
    //       <p>how are you?</p>
    //     </i18n>
    //     -> passed values: Array (included VNode):
    //        `[VNode{ tag: 'p', text: 'kazupon', ...}, VNode{ tag: 'p', text: 'how are you?', ...}]`
    //
    // @return {Array<any>}
    //   interpolated values. you need to return the following:
    //   - array of string, when is using `$t` or `$tc`.
    //   - array included VNode object, when is using `i18n` functional component.
    //
    interpolate(message: string, values: Record<string, unknown>): (string | any[])[] {
        const customFormatters = {
            projectDuration: (v) => {
                return _.capitalize(moment.duration(v / 1000, "seconds").humanize());
            },
        };
        const mf = new MessageFormat("en", { customFormatters });
        const compiled = mf.compile(message);
        return [compiled(values)];
    }
}

Vue.use(VueI18n);

interface ModuleLocales {
    name: string;
    sensors: Record<string, string>;
}

function loadLocaleMessages(): LocaleMessages {
    const locales = require.context("./locales", true, /[A-Za-z0-9-_,\s]+\.json$/i);
    const messages: LocaleMessages = {};

    locales.keys().forEach((key) => {
        const matched = key.match(/([A-Za-z0-9-_]+)\./i);
        if (matched && matched.length > 1) {
            const locale = matched[1];
            if (messages[locale]) {
                Object.assign(messages[locale], locales(key));
            } else {
                messages[locale] = locales(key);
            }
        }
    });

    const keys = _((messages.en.modules as unknown) as Record<string, ModuleLocales>)
        .map((m, moduleKey) => {
            return _(m.sensors)
                .map((sensorName, sensorKey) => {
                    const normalizedKey = sensorKey
                        .split(".")
                        .map((p) =>
                            _.camelCase(p)
                                .replace("10M", "10m")
                                .replace("2M", "2m")
                        )
                        .join(".");

                    if (moduleKey.indexOf("wh.") == 0) {
                        const fullKey = [moduleKey, normalizedKey].join(".");
                        return [fullKey, sensorName];
                    }

                    const fullKey = ["fk", moduleKey, normalizedKey].join(".");
                    return [fullKey, sensorName];
                })
                .value();
        })
        .flatten()
        .fromPairs()
        .value();

    const moduleKeys = _((messages.en.modules as unknown) as Record<string, ModuleLocales>)
        .map((m, moduleKey) => {
            if (moduleKey.startsWith("wh.")) {
                // HACK
                return [moduleKey, m.name];
            }
            return ["fk." + moduleKey, m.name];
        })
        .fromPairs()
        .value();

    Object.assign(messages.en, keys);
    Object.assign(messages.en, moduleKeys);

    return messages;
}

export default new VueI18n({
    locale: process.env.VUE_APP_I18N_LOCALE || "en",
    fallbackLocale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || "en",
    messages: loadLocaleMessages(),
    formatter: new MessageFormatFormatter(),
});
