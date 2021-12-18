import 'package:angular/angular.dart';

import '../position_component/position_component.dart';
import '../editor_component/editor_component.dart';

@Component(
  selector: 'lcnc-home',
  styleUrls: ['home_component.css'],
  templateUrl: 'home_component.html',
  directives: [PositionComponent, EditorComponent],
)
class HomeComponent {
  // Nothing here yet. All logic is in TodoListComponent.
}
