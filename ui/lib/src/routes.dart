import 'package:angular_router/angular_router.dart';

import 'route_paths.dart';
import 'home_component/home_component.template.dart' as home_template;
import 'probing_component/probing_component.template.dart' as probing_template;

export 'route_paths.dart';

class Routes {
  static final home = RouteDefinition(
    useAsDefault: true,
    routePath: RoutePaths.home,
    component: home_template.HomeComponentNgFactory,
  );

  static final probing = RouteDefinition(
    routePath: RoutePaths.probing,
    component: probing_template.ProbingComponentNgFactory,
  );

  static final all = <RouteDefinition>[
    home,
    probing,
  ];
}
