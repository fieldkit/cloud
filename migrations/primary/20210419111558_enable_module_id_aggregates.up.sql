ALTER TABLE fieldkit.aggregated_10s RENAME TO aggregated_old_10s;
ALTER TABLE fieldkit.aggregated_1m RENAME TO aggregated_old_1m;
ALTER TABLE fieldkit.aggregated_10m RENAME TO aggregated_old_10m;
ALTER TABLE fieldkit.aggregated_30m RENAME TO aggregated_old_30m;
ALTER TABLE fieldkit.aggregated_1h RENAME TO aggregated_old_1h;
ALTER TABLE fieldkit.aggregated_6h RENAME TO aggregated_old_6h;
ALTER TABLE fieldkit.aggregated_12h RENAME TO aggregated_old_12h;
ALTER TABLE fieldkit.aggregated_24h RENAME TO aggregated_old_24h;

ALTER TABLE fieldkit.aggregated_bymod_10s RENAME TO aggregated_10s;
ALTER TABLE fieldkit.aggregated_bymod_1m RENAME TO aggregated_1m;
ALTER TABLE fieldkit.aggregated_bymod_10m RENAME TO aggregated_10m;
ALTER TABLE fieldkit.aggregated_bymod_30m RENAME TO aggregated_30m;
ALTER TABLE fieldkit.aggregated_bymod_1h RENAME TO aggregated_1h;
ALTER TABLE fieldkit.aggregated_bymod_6h RENAME TO aggregated_6h;
ALTER TABLE fieldkit.aggregated_bymod_12h RENAME TO aggregated_12h;
ALTER TABLE fieldkit.aggregated_bymod_24h RENAME TO aggregated_24h;
