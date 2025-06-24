box.cfg {
  listen = 3301
}

box.once("schema", function()
  box.schema.space.create("kv")
  box.space.kv:format({
    { name = "key",   type = "string" },
    { name = "value", type = "string" }
  })
  box.space.kv:create_index("primary", {
    parts = { "key" }
  })
end)
