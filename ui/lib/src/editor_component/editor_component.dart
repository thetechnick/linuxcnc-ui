import 'dart:html';
import 'dart:js';

import 'package:js/js.dart';
import 'package:angular/angular.dart';
import 'package:amdjs/amdjs.dart';

@Component(
  selector: 'lcnc-editor',
  // styleUrls: ['editor_component.css'],
  styles: [
    '''
      :host {
          display: block;
          height: 200px;
      }
      .editor-container {
          width: 100%;
          height: 98%;
      }
  '''
  ],
  template: '''<div class="editor-container" #editorContainer></div>''',
  directives: [],
)
class EditorComponent implements AfterViewInit {
  @ViewChild('editorContainer')
  HtmlElement test;

  @override
  Future<void> ngAfterViewInit() async {
    const url = 'external/monaco-editor-0.31.1/package/min/vs';

    // Disable DDC onResourceLoad listener while loading,
    // to not end up in an infinite loop.
    //
    // See:
    // https: //github.com/dart-lang/webdev/blob/0274a13ffd50a77a4808c294d843ec1f7dcd06cd/dwds/lib/src/loaders/require.dart#L176-L209
    context.callMethod('eval', [
      '__ddcOnResourceLoad = requirejs.onResourceLoad; requirejs.onResourceLoad = null;'
    ]);

    await AMDJS.require(
      'vs',
      jsSubPath: 'editor/editor.main.js',
      // jsLocation: '/editor/editor.main',
      jsLocation: url,
      // globalJSVariableName: 'monaco',
      addScriptTagInsideBody: true,
    );

    // Re-enable DDC onResourceLoad listener
    context.callMethod(
        'eval', ['requirejs.onResourceLoad = __ddcOnResourceLoad;']);

    var editor = create(
        test,
        MonacoOptions(
          value: ['function x() {', '\tconsole.log("Hello world!");', '}']
              .join('\n'),
          language: 'javascript',
          theme: 'vs-dark',
        ));

    window.onResize.listen((event) {
      print('resize!');
      editor.layout();
    });
  }
}

@JS('monaco.editor.create')
external Monaco create(HtmlElement domElement, MonacoOptions opts);

@JS()
class Monaco {
  external void layout();
}

@JS()
@anonymous
class MonacoOptions {
  external String get language;
  external String get value;
  external String get theme;

  external factory MonacoOptions({String language, String value, String theme});
}
