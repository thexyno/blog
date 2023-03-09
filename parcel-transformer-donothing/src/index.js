const {Transformer} = require('@parcel/plugin');

exports.default = new Transformer({
        async transform({asset}) {
                return [asset];
        }
});
