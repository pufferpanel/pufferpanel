use rhai::{Engine, Scope};
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub fn resolve_if(script: &str, data: &JsValue) -> Result<bool, String> {
    let engine = Engine::new_raw();

    let mut scope = Scope::new();

    let keys = match js_sys::Reflect::own_keys(data) {
        Ok(res) => res,
        Err(e) => return Err(e.as_string().unwrap())
    };

    for key in keys {
        let value = match js_sys::Reflect::get(data, &key) {
            Ok(res) => res,
            Err(e) => return Err(e.as_string().unwrap())
        };

        let k = key.as_string().unwrap();
        if value.as_bool().is_some() {
            scope.push(k, value.as_bool().unwrap());
        } else if value.as_f64().is_some() {
            scope.push(k, value.as_f64().unwrap());
        } else if value.as_string().is_some() {
            scope.push(k, value.as_string().unwrap());
        }
    }

    let result = match engine.eval_expression_with_scope::<bool>(&mut scope, script) {
        Ok(res) => res,
        Err(e) => return Err(e.to_string())
    };
    Ok(result)
}
