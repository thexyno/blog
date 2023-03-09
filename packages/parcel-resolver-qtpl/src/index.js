const {Resolver} = require('@parcel/plugin');

exports.default = new Resolver({
        async resolve(x) {
                if(!x || !x.dependency | !x.dependency.sourcePath) return null;
                if(x.dependency.sourcePath.endsWith(".qtpl") && !(x.specifier.endsWith(".css") || x.specifier.endsWith(".js"))) {
                        return { isExcluded: true };
                }
                return null;
        }
});
