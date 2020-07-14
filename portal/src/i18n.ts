import _ from "lodash";
import Vue from "vue";
import VueI18n, { LocaleMessages } from "vue-i18n";

Vue.use(VueI18n);

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

    const keys = _(messages.en.modules)
        .map((module, moduleKey) => {
            return _(module.sensors)
                .map((sensorName, sensorKey) => {
                    const normalizedKey = sensorKey
                        .split(".")
                        .map((p) =>
                            _.camelCase(p)
                                .replace("10M", "10m")
                                .replace("2M", "2m")
                        )
                        .join(".");
                    const fullKey = ["fk", moduleKey, normalizedKey].join(".");
                    return [fullKey, sensorName];
                })
                .value();
        })
        .flatten()
        .fromPairs()
        .value();

    Object.assign(messages.en, keys);

    return messages;
}

export default new VueI18n({
    locale: process.env.VUE_APP_I18N_LOCALE || "en",
    fallbackLocale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || "en",
    messages: loadLocaleMessages(),
});
