let apply = (effects: Model.effects) => Js.log(("apply", effects));

let run = (config: Model.configuration) => {
  let initial = Model.init(config);
  let model = ref(initial);
  let process = (event: Model.event) => {
    let (next, effects) = Sync.update(model^, event);
    model := next;
    apply(effects);
  };
  Js.Global.setInterval(() => process(Tick), 1);
};
