const { Resolver } = require('@parcel/plugin');

exports.default = new Resolver({
  async resolve(x) {
    if (!x.dependency | !x.dependency.sourcePath) return null;
    // make sure only css and js files are included from qtpl files
    if (x.dependency.sourcePath.endsWith(".qtpl") &&
      !(x.specifier.endsWith(".css") || x.specifier.endsWith(".js"))) {
      return { isExcluded: true };
    }
    return null;
  }
});
