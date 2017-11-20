package modules

import "rider"

type SecureConfig struct {

	//0: 禁止XSS过滤。
	//1: 启用XSS过滤（通常浏览器是默认的）。 如果检测到跨站脚本攻击，浏览器将清除页面（删除不安全的部分）。
	//1;mode=block: 启用XSS过滤。 如果检测到攻击，浏览器将不会清除页面，而是阻止页面加载。
	//1; report=<reporting-URI>  (Chromium only): 启用XSS过滤。 如果检测到跨站脚本攻击，浏览器将清除页面并使用CSP report-uri指令的功能发送违规报告。
	XXSSProtection string `header:"X-XSS-Protection"`  //"1; mode=block"

	//阻止点击挟持攻击
	//DENY // 拒绝任何域加载
	//SAMEORIGIN // 允许同源域下加载
	//ALLOW-FROM // 可以定义允许frame加载的页面地址
	XFrameOptions string `header:"x_frame_options"`   //SAMEORIGIN

	//HTTP转发HTTPS是302重定向
	//用户在浏览器输入www.baidu.com；然后重定向到https://www.baidu.com；在这个过程中，这个跳转是可以被劫持的，叫做http劫持，特别用公共路由器
	//服务器端配置HSTS，减少302跳转，其实HSTS的最大作用是防止302 HTTP劫持。HSTS的缺点是浏览器支持率不高，另外配置HSTS后HTTPS很难实时降级成HTTP。
	//详细可以查看http://www.jianshu.com/p/caa80c7ad45c
	//可以配置三个参数max-age；includeSubDomains；preload；
	//max-age是必选参数，是一个以秒为单位的数值，它代表着HSTS Header的过期时间，通常设置为1年，即31536000秒。
	//includeSubDomains是可选参数，如果包含它，则意味着当前域名及其子域名均开启HSTS保护。
	//preload参数可以向https://hstspreload.org官网申请
	//配置案例  Strict-Transport-Security: max-age=300; includeSubDomains
	StrictTransportSecurity string `header:"Strict-Transport-Security"`

	//这就意味着脚本文件只能来自content_security_policy指定的范围
	//via:  script-src self；表明js脚本文件只允许网站自己，但是加了这个响应头，所有的js脚本都得通过<script src="xxx.js"></script>引入，
	//非本站脚本会被阻止，页面内部的<script>xxx</script>的脚本会忽略
	//其他配置还有image-src; style-src...
	ContentSecurityPolicy string `header:"content_security_policy"`

	//X-Content-Type-Options 响应首部相当于一个提示标志，被服务器用来提示客户端一定要遵循在 Content-Type 首部中对  MIME 类型 的设定，而不能对其进行修改。这就禁用了客户端的 MIME 类型嗅探行为，换句话说，也就是意味着网站管理员确定自己的设置没有问题。
	XContentTypeOptions string `header:"X-Content-Type-Options"`  //nosniff
}


var DefaultSecureConfig *SecureConfig = &SecureConfig{
	XXSSProtection: "1; mode=block",
	XFrameOptions: "SAMEORIGIN",
	StrictTransportSecurity: "",
	ContentSecurityPolicy: "",
	XContentTypeOptions: "nosniff",
}


func SecureHeader() rider.HandlerFunc {
	return func(c rider.Context) {
		if DefaultSecureConfig.XXSSProtection != "" {
			c.SetHeader(rider.HeaderXXSSProtection, DefaultSecureConfig.XXSSProtection)
		}
		if DefaultSecureConfig.XFrameOptions != "" {
			c.SetHeader(rider.HeaderXFrameOptions, DefaultSecureConfig.XFrameOptions)
		}
		if DefaultSecureConfig.StrictTransportSecurity != "" {
			c.SetHeader(rider.HeaderStrictTransportSecurity, DefaultSecureConfig.StrictTransportSecurity)
		}
		if DefaultSecureConfig.ContentSecurityPolicy != "" {
			c.SetHeader(rider.HeaderContentSecurityPolicy, DefaultSecureConfig.ContentSecurityPolicy)
		}
		if DefaultSecureConfig.XContentTypeOptions != "" {
			c.SetHeader(rider.HeaderXContentTypeOptions, DefaultSecureConfig.XContentTypeOptions)
		}
		c.Next()
	}
}