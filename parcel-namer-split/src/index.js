const {Namer} = require('@parcel/plugin');
const path = require('node:path');

exports.default = new Namer({
        name({bundle}) {
                if(bundle.type != "qtpl") {
                   let filePath = bundle.getMainEntry().filePath;
                   let bn = path.basename(filePath).split(".")
                   let hr = bundle.needsStableName ? "." : `${bundle.hashReference}.`
                   return `../statics/dist/${bn[0]}.${hr}${bn.slice(1).join("")}`;
                }
                return null;

        }
});
