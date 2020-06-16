const {init, update} = require('./Model.bs.js')
const {tick} = require('./Event.bs.js')

const apply = (effects) => {
  console.log(effects)
}

const run = () => {
  let config = {}
  let model = init(config)
  let process = (event) => {
    [model, effects] = update(model, event)
    apply(effects)
  }
  setInterval(() => process(tick), 1)
}
