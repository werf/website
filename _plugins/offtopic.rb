module Jekyll
  module Offtopic
    class OfftopicTag < Liquid::Block
      def initialize(tag_name, params_as_string, tokens)
        super
        @params_as_string = params_as_string
      end

      def render(context)
        result = ""

        begin
          content = super

          @unnamed_params, @named_params = Utils.parse_params(context, @params_as_string)
          Utils.validate_params(@unnamed_params, @named_params, {
            named: [
              { name: "title" },
              { name: "compact", regex: /^(true)|(false)$/ }
            ]
          })

          unless @named_params.key?("title")
            @named_params["title"] = "Подробности"
          end

          unless @named_params.key?("compact")
            @named_params["compact"] = false
          end

          site_config = context.registers[:site].config
          rendered_content = Jekyll::Converters::Markdown::KramdownParser.new(site_config).convert(content)

          if @named_params["compact"]
            div_details_class = "details details__compact"
          else
            div_details_class = "details"
          end

          result = %Q(
<div class="#{div_details_class}">
<a href="javascript:void(0)" class="details__summary">#{@named_params["title"]}</a>
<div class="details__content" markdown="1">
#{rendered_content}
</div>
</div>
)
        rescue => e
          Jekyll.logger.abort_with("[offtopic] FATAL:", e.message)
        end

        result
      end
    end
  end
end

Liquid::Template.register_tag('offtopic', Jekyll::Offtopic::OfftopicTag)
