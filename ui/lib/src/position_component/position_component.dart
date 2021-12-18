import 'package:angular/angular.dart';

@Component(
  selector: 'lcnc-position',
  styleUrls: ['position_component.css'],
  templateUrl: 'position_component.html',
  directives: [],
  pipes: [commonPipes],
)
class PositionComponent {
  // Nothing here yet. All logic is in TodoListComponent.
  @Input()
  String title;

  @Input()
  double x;
  @Input()
  double y, z, a, b, c, u, v, w;
}
