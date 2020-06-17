let apply = effects => Js.log(("apply", effects));

let run = (config: Model.configuration) => {
  let initial = Model.init(config);
  let model = ref(initial);
  let process = (event: Event.event) => {
    Js.log(("process", event));
    let (next, effects) = Model.update(model^, event);
    model := next;
    apply(effects);
  };
  Js.Global.setInterval(() => process(Tick), 1);
};
